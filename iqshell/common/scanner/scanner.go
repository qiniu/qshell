package scanner

import (
	"bufio"
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

type Info struct {
	StdInEnable bool // true: InputFile 未设置时使用 stdin
	InputFile   string
}

type Scanner interface {
	LineCount() int64
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
		s.lineCount, _ = utils.FileLineCounts(info.InputFile)
		s.file = f
		s.scanner = bufio.NewScanner(f)
		log.InfoF("read data from file:%s", info.InputFile)
	} else if info.StdInEnable {
		s.scanner = bufio.NewScanner(os.Stdin)
		log.Info("read data from stdin, you can end input with ctrl + D and cancel by ctrl + C")
	} else {
		return nil, errors.New("no scanner source")
	}
	return s, nil
}

type lineScanner struct {
	lineCount int64
	file      *os.File
	scanner   *bufio.Scanner
}

func (b *lineScanner) LineCount() int64 {
	return b.lineCount
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
