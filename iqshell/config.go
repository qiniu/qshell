package iqshell

import (
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

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
