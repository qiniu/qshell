package upload

import (
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"sync/atomic"
)

type resumeProgress struct {
	uploadSize int64
	fileSize   int64
	pgr        progress.Progress
}

func newResumeProgress(p progress.Progress, fileSize int64) *resumeProgress {
	return &resumeProgress{
		uploadSize: 0,
		fileSize:   fileSize,
		pgr:        p,
	}
}

func (p *resumeProgress) start() {
	if p.pgr != nil {
		p.pgr.Start()
	}
}

func (p *resumeProgress) completeSendBlock(blockSize int64) {
	if p.pgr != nil {
		atomic.AddInt64(&p.uploadSize, blockSize)
		p.pgr.Progress(p.fileSize, p.uploadSize)
	}
}

func (p *resumeProgress) end() {
	if p.pgr != nil {
		p.pgr.End()
	}
}
