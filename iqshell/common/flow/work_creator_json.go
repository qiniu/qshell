package flow

import (
	"encoding/json"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type jsonWorkCreator struct {
	BlankQuotedWorkCreatFunc func() Work
}

func (w *jsonWorkCreator) Create(info string) (Work, *data.CodeError) {
	work := w.BlankQuotedWorkCreatFunc()
	err := json.Unmarshal([]byte(info), work)
	return work, data.ConvertError(err)
}

func NewJsonWorkCreator(blankQuotedWorkCreatFunc func() Work) WorkCreator {
	return &jsonWorkCreator{
		BlankQuotedWorkCreatFunc: blankQuotedWorkCreatFunc,
	}
}
