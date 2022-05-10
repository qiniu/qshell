package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Worker interface {

	// DoWork 处理工作
	// @Description: recordList 长度需和 workInfos 长度想等
	// @param workInfos 工作列表
	// @return recordList 工作记录列表
	// @return err 工作错误信息
	DoWork(workInfos []*WorkInfo) (recordList []*WorkRecord, err *data.CodeError)
}

func NewWorker(doFunc func(workInfos []*WorkInfo) ([]*WorkRecord, *data.CodeError)) Worker {
	return &workerStruct{
		DoFunc: doFunc,
	}
}

func NewSimpleWorker(doFunc func(workInfo *WorkInfo) (Result, *data.CodeError)) Worker {
	return &workerStruct{
		SimpleDoFunc: doFunc,
	}
}

type workerStruct struct {
	SimpleDoFunc func(workInfo *WorkInfo) (Result, *data.CodeError)
	DoFunc       func(workInfos []*WorkInfo) ([]*WorkRecord, *data.CodeError)
}

func (w *workerStruct) DoWork(workInfoList []*WorkInfo) ([]*WorkRecord, *data.CodeError) {
	if w == nil {
		return nil, alert.Error("worker: no worker", "")
	}

	if w.DoFunc != nil {
		return w.DoFunc(workInfoList)
	} else if w.SimpleDoFunc != nil {
		recordList := make([]*WorkRecord, 0, len(workInfoList))
		for _, workInfo := range workInfoList {
			record := &WorkRecord{
				WorkInfo: workInfo,
			}
			record.Result, record.Err = w.SimpleDoFunc(workInfo)
			recordList = append(recordList, record)
		}
		return recordList, nil
	} else {
		return nil, alert.Error("worker: no worker func", "")
	}
}
