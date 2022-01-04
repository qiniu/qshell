package batch

import (
	"encoding/base64"
	"fmt"
)

var (
	defaultOperationCountPerRequest = 1000
)

type OperationCondition struct {
	PutTime string
}

func OperationConditionURI(condition OperationCondition) string {
	cond := ""
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
}
