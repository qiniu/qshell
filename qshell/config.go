package qshell

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// config files, priority high to low

var (
	//dir to store some cached files for qshell, like ak, sk
	QShellRootPath string
	AccountFname   string
	AccountDBPath  string

	QShellConfigFiles [2]string
	UpHost            = "upload.qiniup.com"
	RsfHost           = storage.DefaultRsfHost
	IoHost            = "iovip.qbox.me"
	RsHost            = storage.DefaultRsHost
	ApiHost           = storage.DefaultAPIHost
)

func init() {
	curUser, gErr := user.Current()
	if gErr != nil {
		fmt.Println("Error: get current user error,", gErr)
		os.Exit(STATUS_HALT)
	}
	QShellRootPath = filepath.Join(curUser.HomeDir, ".qshell")
	AccountFname = filepath.Join(QShellRootPath, "account.json")
	AccountDBPath = filepath.Join(QShellRootPath, "account.db")
	QShellConfigFiles = [2]string{filepath.Join("/etc/", "qshellrc"), filepath.Join(QShellRootPath, ".qshellrc")}
}

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)

const (
	STATUS_OK = iota
	//process error
	STATUS_ERROR
	//local error
	STATUS_HALT
)

func parseConfigFile(filePath string) (err error) {
	hostFp, openErr := os.Open(filePath)
	if openErr != nil {
		return
	}
	scanner := bufio.NewScanner(hostFp)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.SplitN(line, "=", 1)
		varName := strings.ToUpper(strings.TrimSpace(splits[0]))
		varValue := strings.ToLower(strings.TrimSpace(splits[1]))

		switch varName {
		case "UP_HOST":
			UpHost = varValue
		case "RS_HOST":
			RsHost = varValue
		case "RSF_HOST":
			RsfHost = varValue
		case "IO_HOST":
			IoHost = varValue
		case "API_HOST":
			ApiHost = varValue
		}
	}
	return
}

func ReadConfigFile() error {

	for _, filePath := range QShellConfigFiles {
		err := parseConfigFile(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
