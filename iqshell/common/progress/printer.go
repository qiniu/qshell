package progress

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"sync"
	"time"
)

const (
	wordsCountPerLine = 80
)

type printer struct {
	mu          sync.Mutex
	title       string
	fileSize    int64
	current     int64
	progressBar *progressbar.ProgressBar
}

func NewPrintProgress(title string) Progress {
	return &printer{
		title: title,
		progressBar: progressbar.NewOptions(0,
			progressbar.OptionFullWidth(),
			progressbar.OptionShowBytes(true),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionThrottle(time.Millisecond*500),
			progressbar.OptionOnCompletion(func() {
				fmt.Printf("\n")
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionSetDescription("[green]"+title+"[reset]"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]-[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			})),
	}
}

var _ Progress = (*printer)(nil)

func (p *printer) Start() {
	_ = p.progressBar.Add(0)
}

func (p *printer) SetFileSize(fileSize int64) {
	p.fileSize = fileSize
	p.progressBar.ChangeMax64(fileSize)
}

func (p *printer) SendSize(newSize int64) {
	p.mu.Lock()
	if p.current+newSize > p.fileSize {
		newSize = p.fileSize - p.current
	}
	_ = p.progressBar.Add(int(newSize))
	p.current += newSize
	p.mu.Unlock()
}

func (p *printer) Write(b []byte) (int, error) {
	if n, e := p.progressBar.Write(b); e != nil {
		return n, e
	} else {
		return n, nil
	}
}

func (p *printer) Progress(current int64) {
	if p.fileSize == 0 {
		return
	}
	if current > p.fileSize {
		current = p.fileSize
	}

	p.mu.Lock()
	newSize := current - p.current
	_ = p.progressBar.Add(int(newSize))
	p.current = current
	p.mu.Unlock()
}

func (p *printer) End() {
	_ = p.progressBar.Add(int(p.fileSize) - int(p.progressBar.State().CurrentBytes))
	_ = p.progressBar.Finish()
}
