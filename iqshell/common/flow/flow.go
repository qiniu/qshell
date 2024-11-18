package flow

import (
	"strings"
	"sync"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/limit"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

type Info struct {
	Force                     bool // 是否强制直接进行 Flow, 不强制需要用户输入验证码验证
	WorkerCount               int  // worker 数量
	MinWorkerCount            int  // 最小 work 数量，当遇到限制错误会减小 work 数，最小 1
	WorkerCountIncreasePeriod int  // WorkerCount 递增的周期，当在 WorkerCountIncreasePeriod 时间内没有遇到限制错误时，会尝试增加 WorkerCount，最小 10s
	StopWhenWorkError         bool // 当某个 work 遇到执行错误是否结束 batch 任务
}

func (i *Info) Check() *data.CodeError {
	if i.WorkerCount < 1 {
		i.WorkerCount = 1
	}

	if i.MinWorkerCount < 1 {
		i.MinWorkerCount = 1
	}

	if i.WorkerCountIncreasePeriod < 10 {
		i.WorkerCountIncreasePeriod = 10
	}

	return nil
}

type Flow struct {
	Info           Info           // flow 的参数信息 【可选】
	WorkProvider   WorkProvider   // work 提供者 【必填】
	WorkerProvider WorkerProvider // worker 提供者 【必填】

	DoWorkInfoListMaxCount int // Worker.DoWork 函数中 works 数组最大长度，默认：250，最小长度为 1
	doWorkInfoListCount    int // Worker.DoWork 函数中 works 数组长度
	DoWorkInfoListMinCount int // Worker.DoWork 函数中 works 数组最小长度，默认：50，最小长度为 1

	Limit         limit.BlockLimit // 速度限制，用于限制
	EventListener EventListener    // work 处理事项监听者 【可选】
	Overseer      Overseer         // work 监工，涉及 work 是否已处理相关的逻辑 【可选】
	Skipper       Skipper          // work 是否跳过相关逻辑 【可选】
	Redo          Redo             // work 是否需要重新做相关逻辑，有些工作虽然已经做过，但下次处理时可能条件发生变化，需要重新处理 【可选】

	mu                sync.Mutex //
	workErrorHappened bool       // 执行中是否出现错误 【内部变量】
}

func (f *Flow) Check() *data.CodeError {
	if err := f.Info.Check(); err != nil {
		return err
	}

	if f.WorkProvider == nil {
		return alert.CannotEmptyError("WorkProvider", "")
	}
	if f.WorkerProvider == nil {
		return alert.CannotEmptyError("WorkerProvider", "")
	}

	if f.DoWorkInfoListMaxCount < 1 {
		f.DoWorkInfoListMaxCount = 1
	}

	if f.DoWorkInfoListMinCount < 1 {
		f.DoWorkInfoListMinCount = 1
	}

	f.doWorkInfoListCount = f.DoWorkInfoListMaxCount

	return nil
}

func (f *Flow) Start() {
	if e := f.Check(); e != nil {
		log.ErrorF("work flow start error:%v", e)
		return
	}

	if !f.Info.Force && !UserCodeVerification() {
		return
	}

	if err := f.notifyFlowWillStart(); err != nil {
		log.ErrorF("Flow start error:%v", err)
		return
	}

	log.Debug("work flow did start")
	workChan := make(chan []*WorkInfo, f.Info.WorkerCount)
	// 生产者
	go func() {
		log.DebugF("work producer start")

		workList := make([]*WorkInfo, 0, f.doWorkInfoListCount)
		for {
			hasMore, workInfo, err := f.WorkProvider.Provide()
			log.DebugF("work producer get work, hasMore:%v, workInfo: %+v, err: %+v", hasMore, workInfo, err)
			if err != nil {
				workInfoData := ""
				if workInfo != nil {
					workInfoData = workInfo.Data
				}
				if err.Code == data.ErrorCodeParamMissing ||
					err.Code == data.ErrorCodeLineHeader {
					log.DebugF("work producer get work, skip:%s because:%s", workInfoData, err)
					f.notifyWorkSkip(workInfo, nil, err)
				} else {
					// 没有读到任何数据
					if workInfo == nil || len(workInfo.Data) == 0 {
						log.ErrorF("work producer get work fail: %s", err)
						break
					}
					log.DebugF("work producer get work fail, error:%s info:%s", err, workInfoData)
					f.notifyWorkFail(workInfo, err)
				}
				continue
			}

			if workInfo == nil || workInfo.Work == nil {
				if !hasMore {
					log.Info("work producer get work completed")
					break
				} else {
					log.Info("work producer get work fail: work in empty")
					continue
				}
			}

			// 检测 work 是否需要过
			if skip, cause := f.shouldWorkSkip(workInfo); skip {
				log.DebugF("work producer get work, skip:%s cause:%s", workInfo.Data, cause)
				f.notifyWorkSkip(workInfo, nil, cause)
				continue
			}

			// 检测 work 是否已经做过
			if hasDone, workRecord := f.getWorkRecordIfHasDone(workInfo); hasDone {
				if shouldRedo, cause := f.shouldWorkRedo(workInfo, workRecord); !shouldRedo {
					if cause == nil {
						cause = data.NewError(data.ErrorCodeAlreadyDone, "already done")
					}
					cause.Code = data.ErrorCodeAlreadyDone
					f.notifyWorkSkip(workInfo, workRecord.Result, cause)
					continue
				} else {
					if cause == nil {
						log.DebugF("work redo, %s", workInfo.Data)
					} else {
						log.DebugF("work redo, %s because:%v", workInfo.Data, cause.Desc)
					}
				}
			}

			// 通知 work 将要执行
			if shouldContinue, e := f.notifyWorkWillDoing(workInfo); !shouldContinue {
				f.notifyWorkSkip(workInfo, nil, e)
				continue
			}

			workList = append(workList, workInfo)
			if len(workList) >= f.doWorkInfoListCount {
				workChan <- workList
				workList = make([]*WorkInfo, 0, f.DoWorkInfoListMaxCount)
			}
		}

		if len(workList) > 0 {
			workChan <- workList
		}

		close(workChan)
		log.DebugF("work producer   end")
	}()

	// 消费者
	wait := &sync.WaitGroup{}
	wait.Add(f.Info.WorkerCount)
	for i := 0; i < f.Info.WorkerCount; i++ {
		time.Sleep(time.Millisecond * time.Duration(50))
		go func(index int) {
			log.DebugF("work consumer %d start", index)
			defer func() {
				wait.Done()
				log.DebugF("work consumer %d   end", index)
			}()

			worker, err := f.WorkerProvider.Provide()
			if err != nil {
				log.ErrorF("Create Worker Error:%v", err)
				return
			}

			for workList := range workChan {
				if workspace.IsCmdInterrupt() {
					break
				}

				workCount := len(workList)
				log.DebugF("work consumer get works, count:%d", workCount)

				_ = f.limitAcquire(workCount)
				// workRecordList 有数据则长度和 workList 长度相同
				workRecordList, workErr := worker.DoWork(workList)
				f.limitRelease(workCount)
				log.DebugF("work consumer handle works, count:%d error:%+v", workCount, workErr)

				if len(workRecordList) == 0 && workErr != nil {
					log.ErrorF("Do Worker Error:%+v", workErr)
					for _, workInfo := range workList {
						f.handleWorkResult(&WorkRecord{
							WorkInfo: workInfo,
							Result:   nil,
							Err:      workErr,
						})
					}
					break
				}

				f.tryChangeWorkGroupCount(workErr)

				hitLimitCount := 0
				hasTooManyFileError := false
				for _, record := range workRecordList {
					if (record.Result == nil || !record.Result.IsValid()) && record.Err == nil {
						record.Err = workErr
					}

					f.handleWorkResult(record)
					if f.isWorkResultHitLimit(record) {
						hitLimitCount += 1
					}

					if !hasTooManyFileError &&
						record.Err != nil &&
						strings.Contains(record.Err.Error(), "too many open files") {
						hasTooManyFileError = true
					}
				}
				f.limitCountDecrease(hitLimitCount)

				if hasTooManyFileError {
					time.Sleep(5 * time.Second)
				}
				// 检测是否需要停止
				if f.workErrorHappened && f.Info.StopWhenWorkError {
					break
				}
			}
		}(i)
	}
	wait.Wait()

	if err := f.notifyFlowWillEnd(); err != nil {
		log.ErrorF("Flow end error:%v", err)
		return
	}

	log.Debug("work flow did end")
}

func (f *Flow) notifyFlowWillStart() *data.CodeError {
	if f.EventListener.FlowWillStartFunc == nil {
		return nil
	}
	return f.EventListener.FlowWillStartFunc(f)
}

func (f *Flow) shouldWorkSkip(work *WorkInfo) (skip bool, cause *data.CodeError) {
	if f.Skipper == nil {
		return false, nil
	}
	return f.Skipper.ShouldSkip(work)
}

func (f *Flow) notifyWorkSkip(work *WorkInfo, result Result, err *data.CodeError) {
	f.EventListener.OnWorkSkip(work, result, err)
}

func (f *Flow) getWorkRecordIfHasDone(work *WorkInfo) (hasDone bool, record *WorkRecord) {
	if f.Overseer == nil {
		return false, nil
	}
	return f.Overseer.GetWorkRecordIfHasDone(work)
}

func (f *Flow) shouldWorkRedo(work *WorkInfo, workRecord *WorkRecord) (shouldRedo bool, cause *data.CodeError) {
	if f.Redo == nil {
		return false, data.NewError(data.ErrorCodeAlreadyDone, workRecord.Err.Error())
	}
	return f.Redo.ShouldRedo(work, workRecord)
}

func (f *Flow) notifyWorkWillDoing(work *WorkInfo) (shouldContinue bool, err *data.CodeError) {
	return f.EventListener.WillWork(work)
}

func (f *Flow) limitAcquire(count int) *data.CodeError {
	if f.Limit == nil {
		return nil
	}
	return f.Limit.Acquire(count)
}

func (f *Flow) isWorkResultHitLimit(workRecord *WorkRecord) bool {
	if f.Limit == nil || workRecord.Err == nil {
		return false
	}

	return workRecord.Err.Code == 573
}

func (f *Flow) limitRelease(count int) {
	if f.Limit == nil {
		return
	}
	f.Limit.Release(count)
}

func (f *Flow) limitCountDecrease(count int) {
	if f.Limit == nil || count <= 0 {
		return
	}

	f.Limit.AddLimitCount(-1 * count)
}

func (f *Flow) tryChangeWorkGroupCount(err *data.CodeError) {
	if err == nil {
		return
	}

	if err.Code != 504 {
		return
	}

	f.mu.Lock()
	f.doWorkInfoListCount -= 10
	if f.doWorkInfoListCount < 1 {
		f.doWorkInfoListCount = 1
	}
	f.mu.Unlock()
}

func (f *Flow) handleWorkResult(workRecord *WorkRecord) {
	if f.Overseer != nil {
		f.Overseer.WorkDone(&WorkRecord{
			WorkInfo: workRecord.WorkInfo,
			Result:   workRecord.Result,
			Err:      workRecord.Err,
		})
	}
	if workRecord.Err != nil {
		f.notifyWorkFail(workRecord.WorkInfo, workRecord.Err)
		f.workErrorHappened = true
	} else {
		f.notifyWorkSuccess(workRecord.WorkInfo, workRecord.Result)
	}
}

func (f *Flow) notifyWorkSuccess(work *WorkInfo, result Result) {
	f.EventListener.OnWorkSuccess(work, result)
}

func (f *Flow) notifyWorkFail(work *WorkInfo, err *data.CodeError) {
	f.EventListener.OnWorkFail(work, err)
}

func (f *Flow) notifyFlowWillEnd() *data.CodeError {
	if f.EventListener.FlowWillEndFunc == nil {
		return nil
	}
	return f.EventListener.FlowWillEndFunc(f)
}
