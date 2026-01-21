package flow

import (
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

const UnknownWorkCount = int64(-1)

type WorkProvider interface {
	WorkTotalCount() int64
	Provide() (hasMore bool, work *WorkInfo, err *data.CodeError)
}

func NewWorkProviderOfFile(filepath string, enableStdin bool, creator WorkCreator) (provider WorkProvider, err *data.CodeError) {
	if len(filepath) > 0 {
		return NewFileWorkProvider(filepath, creator)
	}

	if enableStdin {
		log.InfoF("input info with stdin, you can end the input with Ctrl-D or cancel the task with Ctrl-C")
		return NewReaderWorkProvider(os.Stdin, creator)
	}

	return nil, alert.CannotEmptyError("FilePath (WorkProviderOfFile)", "")
}
