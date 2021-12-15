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
	Close() error
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

func Empty() Export {
	return &export{}
}

type export struct {
	file   *os.File
	lock   sync.RWMutex
	writer *bufio.Writer
}

var _ Export = (*export)(nil)

func (e *export) Close() error {
	if e == nil || e.file == nil {
		return nil
	}
	return e.file.Close()
}

func (e *export) Export(a ...interface{}) {
	e.export(fmt.Sprint(a...))
}

func (e *export) ExportF(format string, a ...interface{}) {
	e.export(fmt.Sprintf(format, a...))
}

func (e *export) export(text string) {
	if e != nil && e.writer != nil {
		e.lock.Lock()
		e.writer.WriteString(text)
		e.writer.Flush()
		e.lock.Unlock()
	}
}
