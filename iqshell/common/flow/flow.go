package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"sync"
)

type Info struct {
	WorkerCount       int  // worker 数量
	StopWhenWorkError bool // 当某个 work 遇到执行错误是否结束 batch 任务
}

func (i *Info) Check() error {
	if i.WorkerCount < 1 {
		i.WorkerCount = 1
	}
	return nil
}

type Flow struct {
	Info              Info
	WorkProvider      WorkProvider
	WorkPacker        *WorkPacker
	WorkerProvider    WorkerProvider
	EventListener     EventListener
	Overseer          Overseer
	Redo              Redo
	workErrorHappened bool // 执行中是否出现错误
}

func (f *Flow) Check() error {
	if err := f.Info.Check(); err != nil {
		return err
	}

	if f.WorkProvider == nil {
		return alert.CannotEmptyError("WorkProvider", "")
	}
	if f.WorkerProvider == nil {
		return alert.CannotEmptyError("WorkerProvider", "")
	}
	return nil
}

func (f *Flow) Start() {
	log.Debug("work flow did start")

	workChan := make(chan Work, f.Info.WorkerCount)
	// 生产者
	go func() {
		log.DebugF("work producer start")
		for {
			hasMore, work, err := f.WorkProvider.Provide()
			if err != nil {
				f.EventListener.OnWorkFail(work, err)
				continue
			}

			if work == nil {
				if !hasMore {
					break
				} else {
					continue
				}
			}

			// 检测 work 是否有问题
			if shouldContinue, err := f.EventListener.WillWork(work); !shouldContinue {
				f.EventListener.OnWorkSkip(work, err)
				continue
			}

			// 检测 work 是否已经做过
			if f.Overseer != nil {
				hasDone, workRecord := f.Overseer.GetWorkRecordIfHasDone(work)
				if hasDone && f.Redo == nil {
					shouldRedo, cause := f.Redo.ShouldRedo(work, workRecord)
					if !shouldRedo {
						f.EventListener.OnWorkSkip(work, cause)
						continue
					}
				}
			}

			// 工作进行打包
			if f.WorkPacker != nil{
				if e := f.WorkPacker.Pack(work); e != nil {
					log.ErrorF("work pack error:%v", e)
					break
				}
				work = f.WorkPacker.GetWorkPackageAndClean(false)
			}

			if work != nil {
				workChan <- work
			}
		}

		if f.WorkPacker != nil {
			if work := f.WorkPacker.GetWorkPackageAndClean(true); work != nil {
				workChan <- work
			}
		}

		close(workChan)
		log.DebugF("work producer   end")
	}()

	// 消费者
	wait := &sync.WaitGroup{}
	wait.Add(f.Info.WorkerCount)
	for i := 0; i < f.Info.WorkerCount; i++ {
		go func(index int) {
			log.DebugF("work consumer %d start", index)
			worker, err := f.WorkerProvider.Provide()
			if err != nil {
				log.ErrorF("Create Worker Error:%v", err)
				return
			}

			for work := range workChan {
				if workspace.IsCmdInterrupt() {
					break
				}

				workResult, workErr := worker.DoWork(work)

				resultHandler := func(record *WorkRecord) {
					if f.Overseer != nil {
						f.Overseer.WorkDone(record)
					}
					if record.Err != nil {
						f.EventListener.OnWorkFail(record.Work, record.Err)
						f.workErrorHappened = true
					} else {
						f.EventListener.OnWorkSuccess(record.Work, record.Result)
					}
				}

				if workPackage, ok := work.(*WorkPackage); ok {
					if _, ok = workResult.(*WorkPackage); !ok {
						log.Error("result type of workPackage work Error: mast be *WorkPackage")
					}
					for _, record := range workPackage.WorkRecords {
						resultHandler(record)
					}
				} else {
					resultHandler(&WorkRecord{
						Work:   work,
						Result: workResult,
						Err:    workErr,
					})
				}

				// 检测是否需要停止
				if f.workErrorHappened && f.Info.StopWhenWorkError {
					break
				}
			}

			wait.Done()
			log.DebugF("work consumer %d   end", index)
		}(i)
	}
	wait.Wait()

	log.Debug("work flow did end")
}
