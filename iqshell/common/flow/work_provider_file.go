package flow

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

func NewFileWorkProvider(filePath string, creator WorkCreator) (WorkProvider, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("FileWorkProvider, open file error:%v", err)
	}

	provider, err := NewReaderWorkProvider(f, creator)
	if err != nil {
		return nil, err
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

func (p *fileWorkProvider) Provide() (hasMore bool, work Work, err error) {
	return p.workProvider.Provide()
}
