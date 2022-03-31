package flow

type Overseer interface {
	WillWork(work Work)
	WorkDone(record *WorkRecord)
	GetWorkRecordIfHasDone(work Work) (hasDone bool, record *WorkRecord)
}

type WorkRecord struct {
	Work   Work   `json:"work"`
	Result Result `json:"result"`
	Err    error  `json:"-"`
}
