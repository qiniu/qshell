package flow

type WorkRecord struct {
	Work   Work   `json:"work"`
	Result Result `json:"result"`
	Err    error  `json:"err"`
}

type Overseer interface {
	WillWork(work Work)
	WorkDone(record *WorkRecord)
	GetWorkRecordIfHasDone(work Work) (hasDone bool, record *WorkRecord)
}
