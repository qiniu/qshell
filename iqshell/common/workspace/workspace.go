package workspace

import (
	"context"
	"errors"
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
	// 用户目录
	userPath = ""
	// 当前账户
	currentAccount *account.Account

	cancelFunc func()
)

// GetConfig 获取之前需要先 Load
func GetConfig() *config.Config {
	return cfg
}

func GetLogConfig() *config.LogSetting {
	if cfg == nil || cfg.Log == nil {
		return nil
	}
	//if data.Empty(cfg.Log.LogFile) {
	//	logFile := filepath.Join(workspacePath, "log", "qshell.log")
	//	cfg.Log.LogFile = data.NewString(logFile)
	//}
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
	if currentAccount == nil {
		 return account.Account{}, errors.New("can't get current user")
	}
	return *currentAccount, nil
}

func GetMac() (mac *qbox.Mac, err error) {
	acc, gErr := GetAccount()
	if gErr != nil {
		err = gErr
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
