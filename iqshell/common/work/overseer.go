package work

import (
	"encoding/json"
	"github.com/qiniu/qshell/v2/iqshell/common/recorder"
)

type Overseer interface {
	HasDone(work Work) bool
	WillWork(work Work)
	WorkDone(work Work, result Result, err *data.CodeError)
}

func NewDBRecordOverseer(dbPath string) (Overseer, *data.CodeError) {
	if r, err := recorder.CreateDBRecorder(dbPath); err != nil {
		return nil, err
	} else {
		return &localDBRecordOverseer{
			Recorder: r,
		}, nil
	}
}

type localDBRecordOverseer struct {
	Recorder recorder.Recorder
}

func (l *localDBRecordOverseer) HasDone(work Work) bool {
	if l == nil || l.Recorder == nil {
		return false
	}
	s := l.getWorkStatus(work)
	return s.Status == workStatusDoing || s.Status == workStatusSuccess || s.Status == workStatusError
}

func (l *localDBRecordOverseer) WillWork(work Work) {
	if l == nil || l.Recorder == nil {
		return
	}
	s := l.getWorkStatus(work)
	s.Status = workStatusDoing
	l.setWorkStatus(work, s)
}

func (l *localDBRecordOverseer) WorkDone(work Work, result Result, err *data.CodeError) {
	if l == nil || l.Recorder == nil {
		return
	}
	s := l.getWorkStatus(work)
	if err != nil {
		s.Status = workStatusError
	} else {
		s.Status = workStatusSuccess
	}
	l.setWorkStatus(work, s)
}

func (l *localDBRecordOverseer) getWorkStatus(work Work) *workStatus {
	if l == nil || l.Recorder == nil {
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

	s, _ := newWorkStatus(value)
	return s
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

func newWorkStatus(data string) (*workStatus, *data.CodeError) {
	s := &workStatus{}
	err := json.Unmarshal([]byte(data), s)
	return s, err
}

type workStatus struct {
	Status int `json:"status"`
}

func (s *workStatus) toData() (string, *data.CodeError) {
	data, err := json.Marshal(s)
	return string(data), err
}
