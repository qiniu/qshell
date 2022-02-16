package docs

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var documentInfo = make(map[string]string)

func addCmdDocumentInfo(cmdName string, document string) {
	documentInfo[cmdName] = document
}

func ShowCmdDocument(name string) {
	document := documentInfo[name]
	if len(document) == 0 {
		fmt.Printf("doesn't document for cmd:%s \n", name)
	} else {
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
}
