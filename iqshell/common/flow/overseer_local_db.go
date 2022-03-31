package flow

import (
	"encoding/json"
	"github.com/qiniu/qshell/v2/iqshell/common/recorder"
)

func NewDBRecordOverseer(dbPath string, blankWorkRecordBuilder func() *WorkRecord) (Overseer, error) {
	if r, err := recorder.CreateDBRecorder(dbPath); err != nil {
		return nil, err
	} else {
		return &localDBRecordOverseer{
			Recorder:               r,
			BlankWorkRecordBuilder: blankWorkRecordBuilder,
		}, nil
	}
}

type localDBRecordOverseer struct {
	Recorder               recorder.Recorder
	BlankWorkRecordBuilder func() *WorkRecord
}

func (l *localDBRecordOverseer) WillWork(work Work) {
	if l == nil || l.Recorder == nil {
		return
	}
	status := &workStatus{
		WorkRecord: &WorkRecord{
			Work: work,
		},
		Status: workStatusDoing,
	}
	l.setWorkStatus(work, status)
}

func (l *localDBRecordOverseer) WorkDone(record *WorkRecord) {
	if l == nil || l.Recorder == nil {
		return
	}
	status := &workStatus{
		WorkRecord: record,
	}
	if record.Err != nil {
		status.Status = workStatusError
	} else {
		status.Status = workStatusSuccess
	}
	l.setWorkStatus(record.Work, status)
}

func (l *localDBRecordOverseer) GetWorkRecordIfHasDone(work Work) (hasDone bool, record *WorkRecord) {
	if l == nil || l.Recorder == nil {
		return false, nil
	}
	if status := l.getWorkStatus(work); status != nil &&
		(status.Status == workStatusSuccess || status.Status == workStatusError) {
		return true, status.WorkRecord
	} else {
		return false, nil
	}
}

func (l *localDBRecordOverseer) getWorkStatus(work Work) *workStatus {
	if l == nil || l.Recorder == nil || l.BlankWorkRecordBuilder == nil {
		return nil
	}

	workId := work.WorkId()
	if len(workId) == 0 {
		return nil
	}

	value, err := l.Recorder.Get(workId)
	if len(value) == 0 || err != nil {
		return nil
	}

	status := &workStatus{
		WorkRecord: l.BlankWorkRecordBuilder(),
		Status:     workStatusPrepare,
	}

	if e := unmarshalWorkStatus(value, status); e != nil {
		return nil
	} else {
		return status
	}
}

func (l *localDBRecordOverseer) setWorkStatus(work Work, status *workStatus) {
	if l == nil || l.Recorder == nil || status == nil {
		return
	}

	workId := work.WorkId()
	if len(workId) == 0 {
		return
	}

	value, err := status.toData()
	if err != nil {
		return
	}

	_ = l.Recorder.Put(workId, value)
}

var (
	workStatusPrepare = 0
	workStatusDoing   = 1
	workStatusSuccess = 2
	workStatusError   = 3
)

func unmarshalWorkStatus(data string, s *workStatus) error {
	return json.Unmarshal([]byte(data), s)
}

type workStatus struct {
	*WorkRecord

	Status int `json:"status"`
}

func (s *workStatus) toData() (string, error) {
	data, err := json.Marshal(s)
	return string(data), err
}
