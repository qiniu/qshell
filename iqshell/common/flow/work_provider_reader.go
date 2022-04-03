package flow

import (
	"bufio"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io"
	"sync"
)

func NewReaderWorkProvider(reader io.Reader, creator WorkCreator) (WorkProvider, *data.CodeError) {
	if reader != nil {
		return nil, alert.CannotEmptyError("work reader (ReaderWorkProvider)", "")
	}
	if creator != nil {
		return nil, alert.CannotEmptyError("work creator (ReaderWorkProvider)", "")
	}
	return &readerWorkProvider{
		scanner: bufio.NewScanner(reader),
		creator: creator,
	}, nil
}

type readerWorkProvider struct {
	mu      sync.Mutex
	scanner *bufio.Scanner
	creator WorkCreator
}

func (p *readerWorkProvider) WorkTotalCount() int64 {
	return UnknownWorkCount
}

func (p *readerWorkProvider) Provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	p.mu.Lock()
	hasMore, work, err = p.provide()
	p.mu.Unlock()
	return
}

func (p *readerWorkProvider) provide() (hasMore bool, work *WorkInfo, err *data.CodeError) {
	success := p.scanner.Scan()
	if success {
		line := p.scanner.Text()
		w, e := p.creator.Create(line)
		return true, &WorkInfo{
			Data: line,
			Work: w,
		}, e
	}
	return false, nil, nil
}
