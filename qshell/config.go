package qshell

import (
	"bufio"
	"fmt"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
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
	QShellConfigFiles = [2]string{filepath.Join(QShellRootPath, ".qshellrc"),
		filepath.Join("/etc", "qshellrc")}
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

const (
	UpHostId = iota
	ApiHostId
	IoHostId
	RsHostId
	RsfHostId
	BUCKET_RS_HOST_ID
	BUCKET_API_HOST_ID
)

var hosts = map[int]string{
	UpHostId:           "",
	ApiHostId:          "",
	IoHostId:           "",
	RsHostId:           "",
	RsfHostId:          "",
	BUCKET_RS_HOST_ID:  "",
	BUCKET_API_HOST_ID: "",
}

func UpHost() string {
	return hosts[UpHostId]
}

func setHost(hostId int, mac *digest.Mac, bucket, host string) {
	if host == "" {
		if hosts[hostId] == "" {
			//get bucket zone info
			bucketInfo, gErr := GetBucketInfo(mac, bucket)
			if gErr != nil {
				fmt.Println("Get bucket region info error,", gErr)
				os.Exit(STATUS_ERROR)
			}

			//set up host
			SetZone(bucketInfo.Region)
		}
	} else {
		hosts[hostId] = host
	}
}

func SetUpHost(mac *digest.Mac, bucket, host string) {
	setHost(UpHostId, mac, bucket, host)
}

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
			if hosts[UpHostId] == "" {
				hosts[UpHostId] = varValue
			}
		case "RS_HOST":
			if hosts[RsHostId] == "" {
				hosts[RsHostId] = varValue
			}
		case "RSF_HOST":
			if hosts[RsfHostId] == "" {
				hosts[RsfHostId] = varValue
			}
		case "IO_HOST":
			if hosts[IoHostId] == "" {
				hosts[IoHostId] = varValue
			}
		case "API_HOST":
			if hosts[ApiHostId] == "" {
				hosts[ApiHostId] = varValue
			}
		case "BUCKET_RS_HOST":
			if hosts[BUCKET_RS_HOST_ID] == "" {
				hosts[BUCKET_RS_HOST_ID] = varValue
			}
		case "BUCKET_API_HOST":
			if hosts[BUCKET_API_HOST_ID] == "" {
				hosts[BUCKET_API_HOST_ID] = varValue
			}
		}
	}
	return
}

func readConfigFile() {

	for _, filePath := range QShellConfigFiles {
		parseConfigFile(filePath)
	}
}
