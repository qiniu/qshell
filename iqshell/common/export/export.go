package export

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Exporter interface {
	Export(a ...interface{})
	ExportF(format string, a ...interface{})
	Close() error
}

func New(file string) (Exporter, error) {
	if len(file) == 0 {
		return empty(), nil
	}

	fileHandler, err := os.Create(file)
	if err != nil {
		err = fmt.Errorf("open file: %s: %v\n", file, err)
		return empty(), err
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

func (e *exporter) Close() error {
	if e == nil || e.file == nil {
		return nil
	}
	return e.file.Close()
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
