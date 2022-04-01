package flow

import (
	"encoding/json"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type jsonWorkCompactor struct {
}

func (j *jsonWorkCompactor) Compact(work Work) (string, *data.CodeError) {
	d, err := json.Marshal(work)
	return string(d), data.ConvertError(err)
}

func NewJsonWorkCompactor() WorkCompactor {
	return &jsonWorkCompactor{}
}
