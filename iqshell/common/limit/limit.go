package limit

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Limit interface {
	Acquire(count int64) *data.CodeError
	Release(count int64)
}
