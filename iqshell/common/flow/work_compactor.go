package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type WorkCompactor interface {
	Compact(work Work) (info string, err *data.CodeError)
}
