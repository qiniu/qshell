package batch

import (
	"encoding/base64"
	"fmt"
)

var (
	defaultOperationCountPerRequest = 1000
)

// OperationCondition
// 参考链接：https://github.com/qbox/product/blob/eb21b8c26f20e967fa51b210d267c5a4d5ca2af7/kodo/rs.md#delete-%E5%88%A0%E9%99%A4%E8%B5%84%E6%BA%90
type OperationCondition struct {
	FileHash string
	FileMime string
	FileSize string
	PutTime  string
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
	ToOperation() (string, error)
}

type OperationResult struct {
	Code     int
	Hash     string
	FSize    int64
	PutTime  int64
	MimeType string
	Type     int
	Error    string
	Parts    []int64
}
