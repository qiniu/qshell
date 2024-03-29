package export

import (
	"bufio"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"sync"
)

type Exporter interface {
	Export(a ...interface{})
	ExportF(format string, a ...interface{})
	Close() *data.CodeError
}

func New(file string) (Exporter, *data.CodeError) {
	if len(file) == 0 {
		return empty(), nil
	}

	fileHandler, err := os.Create(file)
	if err != nil {
		err = data.NewEmptyError().AppendDescF("open file: %s: %v\n", file, err)
		return empty(), data.NewEmptyError().AppendDesc("open file:" + file).AppendError(err)
	}

	return &exporter{
		file:   fileHandler,
		lock:   sync.RWMutex{},
		writer: bufio.NewWriter(fileHandler),
	}, nil
}

func empty() Exporter {
	return &exporter{}
}

type exporter struct {
	file   *os.File
	lock   sync.RWMutex
	writer *bufio.Writer
}

var _ Exporter = (*exporter)(nil)

func (e *exporter) Close() *data.CodeError {
	if e == nil || e.file == nil {
		return nil
	}

	if err := e.file.Close(); err != nil {
		return data.NewEmptyError().AppendError(err)
	}
	return nil
}

func (e *exporter) Export(a ...interface{}) {
	e.export(fmt.Sprint(a...))
}

func (e *exporter) ExportF(format string, a ...interface{}) {
	e.export(fmt.Sprintf(format, a...))
}

func (e *exporter) export(text string) {
	if e != nil && e.writer != nil {
		e.lock.Lock()
		e.writer.WriteString(text + "\n")
		e.writer.Flush()
		e.lock.Unlock()
	}
}
