package object

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type StatusApiInfo struct {
	Bucket   string
	Key      string
	NeedPart bool
}

func (s StatusApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeStatus|%s|%s|%t", s.Bucket, s.Key, s.NeedPart)
}

type StatusResult struct {
	batch.OperationResult

	// 归档存储文件的解冻状态，uint32 类型，2表示解冻完成，1表示解冻中；归档文件冻结时，不返回该字段。
	RestoreStatus int `json:"restoreStatus"`
	// 文件状态，uint32 类型。1 表示禁用；只有禁用状态的文件才会返回该字段。
	Status int `json:"status"`
	// 文件 md5 值
	MD5 string `json:"md5"`
	// 文件过期删除日期，int64 类型，Unix 时间戳格式
	Expiration int64 `json:"expiration"`
}

func (s StatusApiInfo) ToOperation() (string, *data.CodeError) {
	if len(s.Bucket) == 0 || len(s.Key) == 0 {
		return "", alert.CannotEmptyError("status operation bucket or key", "")
	}
	return storage.URIStat(s.Bucket, s.Key), nil
}

func Status(info StatusApiInfo) (res StatusResult, err *data.CodeError) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		err = data.NewEmptyError().AppendDescF("status object [%s:%s] error:%v", info.Bucket, info.Key, err.Error())
		return
	}

	reqHost, reqErr := bucketManager.RsReqHost(info.Bucket)
	if reqErr != nil {
		err = data.ConvertError(reqErr)
		return
	}

	needParts := "false"
	if info.NeedPart {
		needParts = "true"
	}
	reqURL := fmt.Sprintf("%s%s?needparts=%s", reqHost, storage.URIStat(info.Bucket, info.Key), needParts)
	cErr := bucketManager.Client.CredentialedCall(context.Background(), bucketManager.Mac, auth.TokenQiniu, &res, "POST", reqURL, nil)
	if cErr != nil {
		err = data.NewEmptyError().AppendDescF("status object [%s:%s] status error:%v", info.Bucket, info.Key, cErr.Error())
		return
	}
	return res, nil
}

// ChangeStatusApiInfo 修改 status
type ChangeStatusApiInfo struct {
	Bucket string
	Key    string
	Status int
}

func (c *ChangeStatusApiInfo) ToOperation() (string, *data.CodeError) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", alert.CannotEmptyError("change status operation bucket or key", "")
	}
	return fmt.Sprintf("/chstatus/%s/status/%d", storage.EncodedEntry(c.Bucket, c.Key), c.Status), nil
}

func (c *ChangeStatusApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeStatus|%s|%s|%d", c.Bucket, c.Key, c.Status)
}

func ChangeStatus(info *ChangeStatusApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
