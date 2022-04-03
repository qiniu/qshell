package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

func NewFileWorkProvider(filePath string, creator WorkCreator) (WorkProvider, *data.CodeError) {
	f, oErr := os.Open(filePath)
	if oErr != nil {
		return nil, data.NewEmptyError().AppendDescF("FileWorkProvider, open file error:%v", oErr)
	}

	provider, rErr := NewReaderWorkProvider(f, creator)
	if rErr != nil {
		return nil, data.ConvertError(rErr)
	}

	workCount, fErr := utils.FileLineCounts(filePath)
	if fErr != nil {
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

func (p *fileWorkProvider) Provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	return p.workProvider.Provide()
}
