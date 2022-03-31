package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"os"
)

const UnknownWorkCount = int64(-1)

type WorkProvider interface {
	WorkTotalCount() int64
	Provide() (hasMore bool, work Work, err error)
}

func NewWorkProviderOfFile(filepath string, enableStdIn bool, creator WorkCreator) (provider WorkProvider, err error) {
	if len(filepath) > 0 {
		return NewFileWorkProvider(filepath, creator)
	}

	if enableStdIn {
		return NewReaderWorkProvider(os.Stdin, creator)
	}

	return nil, alert.CannotEmptyError("FilePath (WorkProviderOfFile)", "")
}
