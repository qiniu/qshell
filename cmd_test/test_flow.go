package cmd

import (
	"fmt"
	"os/exec"
	"testing"
)

type lineWriter struct {
	WriteStringFunc func(line string)
}

func (w *lineWriter) Write(p []byte) (n int, err error) {
	line := string(p)
	fmt.Printf("line:%s", line)

	if w.WriteStringFunc != nil {
		w.WriteStringFunc(line)
	}
	return len(p), nil
}

type testFlow struct {
	args          []string
	resultHandler func(line string)
	errorHandler  func(line string)
}

func NewTestFlow(args ...string) *testFlow {
	return &testFlow{
		args: args,
	}
}

func (t *testFlow) ResultHandler(handler func(line string)) *testFlow {
	t.resultHandler = handler
	return t
}

func (t *testFlow) ErrorHandler(handler func(line string)) *testFlow {
	t.errorHandler = handler
	return t
}

func (t *testFlow) Run() {
	fmt.Println("")
	fmt.Println("========== CMD Start:", t.args, "==========")

	cmd := exec.Command("qshell", t.args...)
	cmd.Stdout = &lineWriter{
		WriteStringFunc: t.resultHandler,
	}

	cmd.Stderr = &lineWriter{
		WriteStringFunc: t.errorHandler,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("==========CMD   End:", t.args, " error:", err.Error(), "==========")
	} else {
		fmt.Println("==========CMD   End:", t.args, "==========")
	}
	fmt.Println("")
}


func defaultTestErrorHandler(t *testing.T) func(line string) {
	return func(line string) {
		t.Fail()
	}
}
