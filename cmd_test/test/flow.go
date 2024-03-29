package test

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"os/exec"
)

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

	var err error
	if IsDebug() {
		err = t.runByDebug()
	} else {
		err = t.runByCommand()
	}

	if err != nil {
		fmt.Println("========== CMD   End:", t.args, " error:", err.Error(), "==========")
	} else {
		fmt.Println("========== CMD   End:", t.args, "==========")
	}

	fmt.Println("")
}

func (t *testFlow) runByCommand() error {
	docs.SetShowMethod(docs.ShowMethodStdOut)
	cmd := exec.Command("qshell", t.args...)
	cmd.Stdout = newLineWriter(t.resultHandler)
	cmd.Stderr = newLineWriter(t.errorHandler)
	return cmd.Run()
}

func (t *testFlow) runByDebug() error {
	docs.SetStdout(newLineWriter(t.resultHandler))
	docs.SetShowMethod(docs.ShowMethodStdOut)
	data.SetStdout(newLineWriter(t.resultHandler))
	data.SetStderr(newLineWriter(t.errorHandler))
	args := []string{"qshell"}
	os.Args = append(args, t.args...)
	cmd.Execute()
	return nil
}
