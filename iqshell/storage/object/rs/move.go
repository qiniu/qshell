package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

type MoveApiInfo struct {
	SourceBucket string
	SourceKey    string
	DestBucket   string
	DestKey      string
	Force        bool
}

func (m MoveApiInfo) ToOperation() (string, error) {
	if len(m.SourceBucket) == 0 || len(m.SourceKey) == 0 || len(m.DestBucket) == 0 || len(m.DestKey) == 0 {
		return "", errors.New(alert.CannotEmpty("move operation bucket or key of source and dest", ""))
	}

	return storage.URIMove(m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey, m.Force), nil
}

type MoveApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newMoveApiResult(ret storage.BatchOpRet) MoveApiResult {
	return MoveApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func Move(info MoveApiInfo) (MoveApiResult, error) {
	ret, err := BatchOne(info)
	if err != nil {
		return MoveApiResult{}, err
	}
	return newMoveApiResult(ret), err
}

func BatchMove(infos []MoveApiInfo) ([]MoveApiResult, error) {
	if len(infos) == 0 {
		return nil, nil
	}

	operations := make([]BatchOperation, len(infos))
	for _, info := range infos {
		operations = append(operations, info)
	}

	result, err := Batch(operations)
	if result == nil || len(result) == 0 {
		return nil, err
	}

	ret := make([]MoveApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newMoveApiResult(item))
	}

	return ret, err
}
