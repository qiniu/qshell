package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type StatusApiInfo struct {
	Bucket string
	Key    string
}

func (s StatusApiInfo) ToOperation() (string, error) {
	if len(s.Bucket) == 0 || len(s.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("status operation bucket or key", ""))
	}
	return storage.URIStat(s.Bucket, s.Key), nil
}

func Status(info StatusApiInfo) (res batch.OperationResult, err error) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		err = fmt.Errorf("status object:[%s|%s] error:%v", info.Bucket, info.Key, err.Error())
		return
	}
	status, err := bucketManager.Stat(info.Bucket, info.Key)
	if err != nil {
		err = fmt.Errorf("status object:[%s|%s] status error:%v", info.Bucket, info.Key, err.Error())
		return
	}
	return batch.OperationResult{
		Code:     200,
		Hash:     status.Hash,
		FSize:    status.Fsize,
		PutTime:  status.PutTime,
		MimeType: status.MimeType,
		Type:     status.Type,
		Error:    "",
		Parts:    status.Parts,
	}, nil
}

// ChangeStatusApiInfo 修改 status
type ChangeStatusApiInfo struct {
	Bucket string
	Key    string
	Status int
}

func (c ChangeStatusApiInfo) ToOperation() (string, error) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("change status operation bucket or key", ""))
	}
	return fmt.Sprintf("/chstatus/%s/status/%c", storage.EncodedEntry(c.Bucket, c.Key), c.Status), nil
}

func ChangeStatus(info ChangeStatusApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}
