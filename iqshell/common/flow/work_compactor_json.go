package flow

import "encoding/json"

type jsonWorkCompactor struct {
}

func (j *jsonWorkCompactor) Compact(work Work) (info string, err error) {
	data, err := json.Marshal(work)
	return string(data), err
}

func NewJsonWorkCompactor() WorkCompactor {
	return &jsonWorkCompactor{}
}
