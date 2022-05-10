package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

const ErrorSeparate = "\tQShellError:"

type WorkCreator interface {
	Create(info string) (work Work, err *data.CodeError)
}
