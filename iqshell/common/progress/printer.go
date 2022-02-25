package progress

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"strings"
	"sync"
)

const (
	wordsCountPerLine = 80
)

type printer struct {
	mu               sync.Mutex
	hasPrintProgress bool
	title            string
	total            int64
	current          int64
}

func NewPrintProgress(title string) Progress {
	return &printer{
		title:            title,
		hasPrintProgress: false,
	}
}

var _ Progress = (*printer)(nil)

func (p *printer) Start() {
	p.printProgress(0, 0)
}

func (p *printer) Progress(total, current int64) {
	if total == 0 {
		return
	}
	p.printProgress(total, current)
}

func (p *printer) End() {
	p.printProgress(p.total, p.total)
}

func (p *printer) printProgress(total, current int64) {
	if current < p.current {
		return
	}
	p.total = total
	p.current = current

	currentString := utils.FormatFileSize(current)
	totalString := "--"
	percentString := "-"
	if total > 0 {
		totalString = utils.FormatFileSize(total)
		percentString = fmt.Sprintf("%.0f", float32(current*100)/float32(total))
	}
	progress := fmt.Sprintf("[%s:%s] %s%%", currentString, totalString, percentString)

	p.mu.Lock()
	if p.hasPrintProgress {
		fmt.Printf("\033[%dA\033[K", 1) // 将光标向上移动一行
	}
	separateStringCount := wordsCountPerLine - len(p.title) - len(progress) - 2
	if separateStringCount < 1 {
		separateStringCount = 1
	}
	separateString := strings.Repeat("-", separateStringCount)
	fmt.Printf("%s %s %s\033[K\n", p.title, separateString, progress) // 输出第二行结果
	p.hasPrintProgress = true
	p.mu.Unlock()
}
