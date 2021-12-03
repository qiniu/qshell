package workspace

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"path/filepath"
)

type Option func(w *workspace)

func Workspace(path string) Option {
	return func(w *workspace) {
		w.workspace = path
	}
}

func UserConfigPath(path string) Option {
	return func(w *workspace) {
		w.userConfigPath = path
	}
}

// 加载工作环境
func Load(options ...Option) (err error) {
	ws := &workspace{}

	// 设置配置
	for _, option := range options {
		option(ws)
	}

	// 检查工作目录
	if len(ws.workspace) == 0 {
		err = errors.New("can't get home dir")
		return
	}
	err = utils.CreateDirIfNotExist(ws.workspace)
	if err != nil {
		return
	}

	// 加载账户
	accountDBPath := filepath.Join(ws.workspace, usersDBName)
	accountPath := filepath.Join(ws.workspace, currentUserFileName)
	oldAccountPath := filepath.Join(ws.workspace, oldUserFileName)
	err = account.Load(account.AccountDBPath(accountDBPath),
		account.AccountPath(accountPath),
		account.OldAccountPath(oldAccountPath))
	if err != nil {
		return
	}

	// 设置配置文件路径
	config.Load(config.UserConfigPath(ws.userConfigPath), config.GlobalConfigPath(ws.globalConfigPath))

	// 加载配置
	cfg.Merge(config.GetUser())
	cfg.Merge(config.GetGlobal())
	cfg.Merge(DefaultConfig())

	currentAccount, err := account.GetAccount()

	if err != nil {
		cfg.Credentials = auth.Credentials{
			AccessKey: currentAccount.AccessKey,
			SecretKey: []byte(currentAccount.SecretKey),
		}
	}

	return
}

type workspace struct {
	workspace        string
	userConfigPath   string
	globalConfigPath string
}
