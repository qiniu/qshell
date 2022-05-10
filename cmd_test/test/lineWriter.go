package test

import (
	"fmt"
	"strings"
)

type lineWriter struct {
	buff            string
	WriteStringFunc func(line string)
}

func newLineWriter(writeStringFunc func(line string)) *lineWriter {
	lw := &lineWriter{
		buff:            "",
		WriteStringFunc: writeStringFunc,
	}
	return lw
}

func (w *lineWriter) Write(p []byte) (n int, err error) {
	w.buff += string(p)

	for len(w.buff) > 0 {
		items := strings.SplitN(w.buff, "\n", 2)
		if len(items) < 2 {
			break
		} else {
			line := items[0]
			fmt.Printf("line:%s\n", line)
			w.buff = items[1]
			if w.WriteStringFunc != nil {
				w.WriteStringFunc(line + "\n")
			}
		}
	}

	return len(p), nil
}

func (w *lineWriter) Close() error {
	return nil
}
