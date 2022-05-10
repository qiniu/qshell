package recorder

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Recorder interface {

	// Get 获取记录
	Get(key string) (value string, err *data.CodeError)

	// Put 添加记录
	Put(key, value string) *data.CodeError

	// Delete 删除记录
	Delete(key string) *data.CodeError
}
