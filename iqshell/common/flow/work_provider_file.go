package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

func NewFileWorkProvider(filePath string, creator WorkCreator) (WorkProvider, *data.CodeError) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, data.NewEmptyError().AppendDescF("FileWorkProvider, open file error:%v", err)
	}

	provider, err := NewReaderWorkProvider(f, creator)
	if err != nil {
		return nil, data.ConvertError(err)
	}

	workCount, err := utils.FileLineCounts(filePath)
	if err != nil {
		workCount = UnknownWorkCount
	}

	return &fileWorkProvider{
		workCount:    workCount,
		workProvider: provider,
	}, nil
}

type fileWorkProvider struct {
	workCount    int64
	workProvider WorkProvider
}

func (p *fileWorkProvider) WorkTotalCount() int64 {
	return p.workCount
}

func (p *fileWorkProvider) Provide() (hasMore bool, work Work, err *data.CodeError) {
	return p.workProvider.Provide()
}
