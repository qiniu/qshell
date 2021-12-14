package export

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"os"
	"sync"
)

type Export interface {
	Export(a ...interface{})
	ExportF(format string, a ...interface{})
}

func New(file string) (Export, error) {
	if len(file) == 0 {
		return nil, errors.New(alert.CannotEmpty("file path", ""))
	}

	fileHandler, err := os.Create(file)
	if err != nil {
		err = fmt.Errorf("open file: %s: %v\n", file, err)
		return nil, err
	}

	return &export{
		file:   fileHandler,
		lock:   sync.RWMutex{},
		writer: bufio.NewWriter(fileHandler),
	}, nil
}

type export struct {
	file   *os.File
	lock   sync.RWMutex
	writer *bufio.Writer
}

var _ Export = (*export)(nil)

func (e *export) Export(a ...interface{}) {
	e.export(fmt.Sprint(a...))
}

func (e *export) ExportF(format string, a ...interface{}) {
	e.export(fmt.Sprintf(format, a...))
}

func (e *export) export(text string) {
	if e.writer != nil {
		e.lock.Lock()
		e.writer.WriteString(text)
		e.writer.Flush()
		e.lock.Unlock()
	}
}
