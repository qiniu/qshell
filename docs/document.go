package docs

import (
	_ "embed"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ShowMethod int

const (
	ShowMethodLess   ShowMethod = 1
	ShowMethodStdOut ShowMethod = 2
)

var sdtout io.Writer = os.Stdout
var showMethod = ShowMethodLess

func SetStdout(o io.Writer) {
	sdtout = o
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
		fmt.Printf("doesn't found document for cmd:%s \n", name)
		return
	}

	if showMethod == ShowMethodStdOut || !utils.IsCmdExist("less") {
		fmt.Fprintln(sdtout, document)
	} else {
		showDocumentByLessCmd(name, document)
	}
}

func showDocumentByLessCmd(name string, document string) {
	errorAlerter := func(err error) {
		fmt.Printf("show document for cmd:%s error:%v", name, err)
	}

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	echoCmd := exec.Command("echo", document)
	echoCmd.Stdout = writer
	echoCmd.Stderr = os.Stderr
	lessCmd := exec.Command("less")
	lessCmd.Stdout = os.Stdout
	lessCmd.Stdin = reader
	lessCmd.Stderr = os.Stderr
	if err := echoCmd.Start(); err != nil {
		errorAlerter(fmt.Errorf("echo start:%v", err))
		return
	}
	if err := lessCmd.Start(); err != nil {
		errorAlerter(fmt.Errorf("less start:%v", err))
		return
	}
	if err := echoCmd.Wait(); err != nil {
		errorAlerter(fmt.Errorf("echo wait:%v", err))
		return
	}
	if err := reader.Close(); err != nil {
		errorAlerter(fmt.Errorf("less reader close:%v", err))
		return
	}
	if err := lessCmd.Wait(); err != nil && !strings.Contains(err.Error(), "read/write on closed pipe") {
		errorAlerter(fmt.Errorf("less wait error%v", err))
		return
	}
}
