package flow

import (
	"bufio"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"io"
	"sync"
)

func NewReaderWorkProvider(reader io.Reader, creator WorkCreator) (WorkProvider, error) {
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

func (p *readerWorkProvider) Provide() (hasMore bool, work Work, err error) {
	p.mu.Lock()
	hasMore, work, err = p.provide()
	p.mu.Unlock()
	return
}

func (p *readerWorkProvider) provide() (hasMore bool, work Work, err error) {
	success := p.scanner.Scan()
	if success {
		hasMore = true
		line := p.scanner.Text()
		work, err = p.creator.Create(line)
	}
	return
}
