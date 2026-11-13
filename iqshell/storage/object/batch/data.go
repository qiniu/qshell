package batch

import (
	"encoding/base64"
	"fmt"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
)

var defaultOperationCountPerRequest = 250

// OperationCondition
// 参考链接：https://github.com/qbox/product/blob/eb21b8c26f20e967fa51b210d267c5a4d5ca2af7/kodo/rs.md#delete-%E5%88%A0%E9%99%A4%E8%B5%84%E6%BA%90
type OperationCondition struct {
	FileHash string `json:"file_hash"`
	FileMime string `json:"file_mime"`
	FileSize string `json:"file_size"`
	PutTime  string `json:"put_time"`
}

func OperationConditionURI(condition OperationCondition) string {
	cond := ""
	if condition.FileHash != "" {
		cond += "hash=" + condition.FileHash
	}
	if condition.FileMime != "" {
		cond += "mime=" + condition.FileMime
	}
	if condition.FileSize != "" {
		cond += "fsize=" + condition.FileSize
	}
	if condition.PutTime != "" {
		cond += "putTime=" + condition.PutTime
	}
	if cond == "" {
		return ""
	}
	return fmt.Sprintf("/cond/%s", base64.URLEncoding.EncodeToString([]byte(cond)))
}

type Operation interface {
	flow.Work

	ToOperation() (string, *data.CodeError)
	GetBucket() string
}

type OperationCreator interface {
	Create(info string) (work Operation, err *data.CodeError)
}

type OperationResult struct {
	Code     int     `json:"code"`
	Hash     string  `json:"hash"`
	FSize    int64   `json:"fsize"`
	PutTime  int64   `json:"putTime"`
	MimeType string  `json:"mimeType"`
	Type     int     `json:"type"`
	EndUser  string  `json:"endUser"`
	Error    string  `json:"error"`
	Parts    []int64 `json:"parts"`

	MD5                            string `json:"md5"`                            // 文件 MD5
	Status                         *int   `json:"status"`                         // 文件禁用状态
	RestoreStatus                  *int   `json:"restoreStatus"`                  // 文件解冻状态
	Expiration                     *int64 `json:"expiration"`                     // 过期删除时间，Unix 秒
	TransitionToIA                 *int64 `json:"transitionToIA"`                 // 转低频时间，Unix 秒
	TransitionToIntelligentTiering *int64 `json:"transitionToIntelligentTiering"` // 转智能分层时间，Unix 秒
	TransitionToArchiveIR          int64  `json:"transitionToArchiveIR"`          // 转归档直读时间，Unix 秒，缺失时为 0
	TransitionToArchive            *int64 `json:"transitionToARCHIVE"`            // 转归档时间，Unix 秒
	TransitionToDeepArchive        *int64 `json:"transitionToDeepArchive"`        // 转深度归档时间，Unix 秒
}

var _ flow.Result = (*OperationResult)(nil)

func (r *OperationResult) IsValid() bool {
	return r.IsSuccess()
}

func (r *OperationResult) IsSuccess() bool {
	if r == nil {
		return false
	}

	return (r.Code == 0 || r.Code == 200) && len(r.Error) == 0
}

func (r *OperationResult) ErrorDescription() string {
	if r == nil || r.IsSuccess() {
		return ""
	}
	return fmt.Sprintf("Code:%d Error:%s", r.Code, r.Error)
}
