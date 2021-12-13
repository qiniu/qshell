package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type StatusApiInfo struct {
	Bucket string
	Key    string
}

var _ batch.BatchOperation = (*StatusApiInfo)(nil)

func (s StatusApiInfo) ToOperation() (string, error) {
	if len(s.Bucket) == 0 || len(s.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("status operation bucket or key", ""))
	}
	return storage.URIStat(s.Bucket, s.Key), nil
}

type StatusApiResult struct {
	Hash     string `json:"hash"`
	FSize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Type     int    `json:"type"`
	Code     int    `json:"code"`
	Error    string `json:"error"`
}

func newStatusApiResult(ret storage.BatchOpRet) StatusApiResult {
	return StatusApiResult{
		Hash:     ret.Data.Hash,
		FSize:    ret.Data.Fsize,
		PutTime:  ret.Data.PutTime,
		MimeType: ret.Data.MimeType,
		Type:     ret.Data.Type,
		Code:     ret.Code,
		Error:    ret.Data.Error,
	}
}

func Status(info StatusApiInfo) (StatusApiResult, error) {
	ret, err := batch.One(info)
	if err != nil {
		return StatusApiResult{}, err
	}
	return newStatusApiResult(ret), err
}

func BatchStatus(infos []StatusApiInfo) ([]StatusApiResult, error) {
	if len(infos) == 0 {
		return nil, nil
	}

	operations := make([]batch.BatchOperation, len(infos))
	for _, info := range infos {
		operations = append(operations, info)
	}

	result, err := batch.Batch(operations)
	if result == nil || len(result) == 0 {
		return nil, err
	}

	ret := make([]StatusApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newStatusApiResult(item))
	}

	return ret, err
}
