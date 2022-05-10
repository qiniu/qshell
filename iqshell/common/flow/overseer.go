package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Overseer interface {
	WillWork(work *WorkInfo)
	WorkDone(record *WorkRecord)
	GetWorkRecordIfHasDone(work *WorkInfo) (hasDone bool, record *WorkRecord)
}

type WorkRecord struct {
	*WorkInfo

	Result Result          `json:"result"`
	Err    *data.CodeError `json:"err"`
}
