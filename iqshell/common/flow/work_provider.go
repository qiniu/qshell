package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
)

const UnknownWorkCount = int64(-1)

type WorkProvider interface {
	WorkTotalCount() int64
	Provide() (hasMore bool, work Work, err *data.CodeError)
}

func NewWorkProviderOfFile(filepath string, enableStdin bool, creator WorkCreator) (provider WorkProvider, err *data.CodeError) {
	if len(filepath) > 0 {
		return NewFileWorkProvider(filepath, creator)
	}

	if enableStdin {
		return NewReaderWorkProvider(os.Stdin, creator)
	}

	return nil, alert.CannotEmptyError("FilePath (WorkProviderOfFile)", "")
}
