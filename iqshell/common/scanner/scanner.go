package scanner

import (
	"bufio"
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

type Info struct {
	StdInEnable bool // true: InputFile 未设置时使用 stdin
	InputFile   string
}

type Scanner interface {
	ScanLine() (line string, success bool)
	Close() error
}

// NewScanner 输入
func NewScanner(info Info) (Scanner, error) {
	s := &lineScanner{}
	if len(info.InputFile) > 0 {
		f, err := os.Open(info.InputFile)
		if err != nil {
			return nil, errors.New("open src dest key map file error")
		}
		s.file = f
		s.scanner = bufio.NewScanner(f)
	} else if info.StdInEnable {
		s.scanner = bufio.NewScanner(os.Stdin)
	} else {
		return nil, errors.New("no scanner source")
	}
	return s, nil
}

type lineScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (b *lineScanner) ScanLine() (line string, success bool) {
	success = b.scanner.Scan()
	if success {
		line = b.scanner.Text()
	}
	log.DebugF("scan line:%s success:%t", line, success)
	return
}

func (b *lineScanner) Close() error {
	if b.file == nil {
		return nil
	}
	return b.file.Close()
}
