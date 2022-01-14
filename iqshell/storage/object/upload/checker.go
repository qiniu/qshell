package upload

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
)

type uploadChecker struct {
	Bucket   string
	Key      string
	Hash     string // 文件 hash
	FileSize int64  // 文件大小
}

func (c *uploadChecker) checkServerExist() (exist bool, err error) {
	status, err := object.Status(object.StatusApiInfo{
		Bucket: c.Bucket,
		Key:    c.Key,
	})
	if err != nil {
		err = errors.New("check server exist: status error:" + err.Error())
		return
	}

	exist = true

	if status.Hash != c.Hash {
		err = fmt.Errorf("check server exist: file exist but not match:[%s:%s]", c.Hash, status.Hash)
		return exist, err
	}

	return exist, nil
}
