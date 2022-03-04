package utils

import "os/exec"

func IsCmdExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
