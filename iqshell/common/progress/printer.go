package progress

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	wordsCountPerLine    = 40
	titleSeparateWord    = "="
	progressSeparateWord = "-"
)

type printer struct {
	mu               sync.Mutex
	hasPrintProgress bool
	title            string
}

func NewPrintProgress(title string) Progress {
	if len(title) == 0{
		title = "action"
	}
	return &printer{
		title:            title,
		hasPrintProgress: false,
	}
}

var _ Progress = (*printer)(nil)

func (p *printer) Start() {
	p.mu.Lock()

	separateStringCount := (wordsCountPerLine - 8 - utf8.RuneCountInString(p.title)) / 2
	separateString := strings.Repeat(titleSeparateWord, separateStringCount)
	fmt.Printf("%s %s Start %s\n", separateString, p.title, separateString)
	p.mu.Unlock()
}

func (p *printer) Progress(total, current int64) {
	if total == 0 {
		return
	}
	p.mu.Lock()
	progress := fmt.Sprintf("[%d:%d]%.2f%%", current, total, float32(current)/float32(total))
	p.printProgress(progress)
	p.mu.Unlock()
}

func (p *printer) End() {
	p.mu.Lock()
	separateStringCount := (wordsCountPerLine - 8 - utf8.RuneCountInString(p.title)) / 2
	separateString := strings.Repeat(titleSeparateWord, separateStringCount)
	fmt.Printf("%s %s   End %s\n", separateString, p.title, separateString)
	p.mu.Unlock()
}

func (p *printer) printProgress(progress string) {
	if p.hasPrintProgress {
		fmt.Printf("\033[%dA\033[K", 1) // 将光标向上移动一行
	}
	separateStringCount := wordsCountPerLine - utf8.RuneCountInString(progress) - 1
	if separateStringCount < 1 {
		separateStringCount = 1
	}
	separateString := strings.Repeat(progressSeparateWord, separateStringCount)
	fmt.Printf("%s %s\033[K\n", separateString, progress) // 输出第二行结果
	p.hasPrintProgress = true
}
