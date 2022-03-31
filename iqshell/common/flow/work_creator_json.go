package flow

import (
	"encoding/json"
)

type jsonWorkCreator struct {
	BlankQuotedWorkCreatFunc func()Work
}

func (w *jsonWorkCreator) Create(info string) (work Work, err error) {
	work = w.BlankQuotedWorkCreatFunc()
	err = json.Unmarshal([]byte(info), work)
	return
}

func NewJsonWorkCreator(blankQuotedWorkCreatFunc func()Work) WorkCreator {
	return &jsonWorkCreator{
		BlankQuotedWorkCreatFunc: blankQuotedWorkCreatFunc,
	}
}

