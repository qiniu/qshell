package test

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
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
	if Debug {
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
	cmd := exec.Command("qshell", t.args...)
	cmd.Stdout = newLineWriter(t.resultHandler)
	cmd.Stderr = newLineWriter(t.errorHandler)
	return cmd.Run()
}

func (t *testFlow) runByDebug() error {
	log.SetStdout(newLineWriter(t.resultHandler))
	log.SetStderr(newLineWriter(t.errorHandler))
	iqshell.SetStdoutColorful(false)
	args := []string{"qshell"}
	os.Args = append(args, t.args...)
	cmd.Execute()
	return nil
}
