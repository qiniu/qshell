package docs

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type ShowMethod int

const (
	ShowMethodLess   ShowMethod = 1
	ShowMethodStdOut ShowMethod = 2
)

var (
	stdout     io.Writer = os.Stdout
	showMethod           = ShowMethodLess
)

func SetStdout(o io.Writer) {
	stdout = o
}

func SetShowMethod(method ShowMethod) {
	showMethod = method
}

var documentInfo = make(map[string]string)

func addCmdDocumentInfo(cmdName string, document string) {
	documentInfo[cmdName] = document
}

func ShowCmdDocument(name string) {
	document := documentInfo[name]
	if len(document) == 0 {
		fmt.Printf("didn't find document for cmd:%s \n", name)
		return
	}

	if showMethod == ShowMethodStdOut || !utils.IsCmdExist("less") {
		fmt.Fprintln(stdout, document)
	} else {
		showDocumentByLessCmd(name, document)
	}
}

func showDocumentByLessCmd(name string, document string) {
	errorAlerter := func(err *data.CodeError) {
		fmt.Printf("show document for cmd:%s error:%v", name, err)
	}

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	lessCmd := exec.Command("less")
	lessCmd.Stdout = stdout
	lessCmd.Stdin = reader
	lessCmd.Stderr = os.Stderr

	if err := lessCmd.Start(); err != nil {
		errorAlerter(data.NewEmptyError().AppendDescF("less start:%v", err))
		return
	}

	if _, err := writer.Write([]byte(document)); err != nil {
		errorAlerter(data.NewEmptyError().AppendDescF("document info write error:%v\n", err))
		return
	}

	if err := reader.Close(); err != nil {
		errorAlerter(data.NewEmptyError().AppendDescF("less reader close:%v", err))
		return
	}

	if err := lessCmd.Wait(); err != nil && !strings.Contains(err.Error(), "read/write on closed pipe") {
		errorAlerter(data.NewEmptyError().AppendDescF("less wait error%v", err))
		return
	}
}
