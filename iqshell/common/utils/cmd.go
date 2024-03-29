package utils

import (
	"bytes"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"os/exec"
)

func IsCmdExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func RunCmd(name string, params []string) (string, *data.CodeError) {
	c := exec.Command(name, params...)

	buff := &bytes.Buffer{}
	c.Stdout = buff
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return "", data.NewEmptyError().AppendDescF("cmd run error:%v", err)
	}
	return buff.String(), nil
}

func CmdExistBySuccess() {
	os.Exit(data.StatusOK)
}

func CmdExistByFail() {
	os.Exit(data.StatusError)
}

func CmdExistByUserCancel() {
	os.Exit(data.StatusUserCancel)
}
