package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Skipper interface {
	ShouldSkip(work Work) (skip bool, cause *data.CodeError)
}
