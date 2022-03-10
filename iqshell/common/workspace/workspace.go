package workspace

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
)

const (
	workspaceName         = ".qshell"
	usersDirName          = "users"
	usersDBName           = "account.db"
	currentUserFileName   = "account.json"
	oldUserFileName       = "old_account.json"
	usersWorkspaceDirName = "workspace"
	taskDirName           = "task"
	taskDBName            = "task.db"
	configFileName        = ".qshell.json"
)

var (
	// config 配置信息
	cfg = &config.Config{}

	// 工作路径
	workspacePath = ""

	cancelFunc func()
)

// GetConfig 获取之前需要先 Load
func GetConfig() *config.Config {
	return cfg
}

func GetLogConfig() *config.LogSetting {
	if cfg == nil {
		return nil
	}
	return cfg.Log
}

func GetStorageConfig() *storage.Config {
	r := cfg.GetRegion()
	if len(cfg.Hosts.GetOneUc()) > 0 {
		storage.SetUcHost(cfg.Hosts.GetOneUc(), cfg.IsUseHttps())
	}

	return &storage.Config{
		UseHTTPS:      cfg.IsUseHttps(),
		Region:        r,
		Zone:          r,
		CentralRsHost: cfg.Hosts.GetOneRs(),
	}
}

func GetWorkspace() string {
	return workspacePath
}

func GetAccount() (account.Account, error) {
	return account.GetAccount()
}

func GetMac() (mac *qbox.Mac, err error) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = fmt.Errorf("get account: %v", gErr)
		return
	}

	mac = qbox.NewMac(acc.AccessKey, acc.SecretKey)
	return
}

func GetContext() context.Context {
	ctx := context.Background()
	ctx, cancelFunc = context.WithCancel(ctx)
	return ctx
}

func Cancel() {
	if cancelFunc != nil {
		cancelFunc()
	}
}
