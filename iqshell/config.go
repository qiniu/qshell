package qshell

import (
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

func init() {
	curUser, gErr := user.Current()
	if gErr != nil {
		fmt.Println("Error: get current user error,", gErr)
		os.Exit(STATUS_HALT)
	}
	rootPath := filepath.Join(curUser.HomeDir, ".qshell")

	viper.SetDefault("path.root_path", filepath.Join(curUser.HomeDir, ".qshell"))
	viper.SetDefault("path.acc_db_path", filepath.Join(rootPath, "account.db"))
	viper.SetDefault("path.acc_path", filepath.Join(rootPath, "account.json"))
	viper.SetDefault("hosts.up_host", "upload.qiniup.com")
	viper.SetDefault("hosts.rs_host", storage.DefaultRsHost)
	viper.SetDefault("hosts.rsf_host", storage.DefaultRsfHost)
	viper.SetDefault("hosts.io_host", "iovip.qbox.me")
	viper.SetDefault("hosts.api_host", storage.DefaultAPIHost)
}

func RootPath() string {
	return viper.GetString("path.root_path")
}

func AccDBPath() string {
	return viper.GetString("path.acc_db_path")
}

func AccPath() string {
	return viper.GetString("path.acc_path")
}

func OldAccPath() string {
	acc_path := viper.GetString("path.acc_path")
	if acc_path == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(acc_path), "old_"+filepath.Base(acc_path))
}

func UpHost() string {
	return viper.GetString("hosts.up_host")
}

func RsHost() string {
	return viper.GetString("hosts.rs_host")
}

func RsfHost() string {
	return viper.GetString("hosts.rsf_host")
}

func IoHost() string {
	return viper.GetString("hosts.io_host")
}

func ApiHost() string {
	return viper.GetString("hosts.api_host")
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

type FetchItem struct {
	RemoteUrl string
	Bucket    string
	Key       string
}
