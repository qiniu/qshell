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

func CmdConfig(cfg *config.Config) Option {
	return func(w *workspace) {
		w.cmdConfig = cfg
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

	// 检查用户路径
	currentAccount, err := account.GetAccount()
	currentAccountDir := ""
	if err == nil {
		accountName := currentAccount.Name
		if len(accountName) == 0 {
			accountName = currentAccount.AccessKey
		}

		currentAccountDir = filepath.Join(ws.workspace, accountName)
		err := utils.CreateDirIfNotExist(currentAccountDir)
		if err != nil {
			return errors.New("create user dir error:" + err.Error())
		}
	}

	// 检查用户配置，用户配置可能被指定，如果未指定则使用用户目录下配置
	if len(ws.userConfigPath) == 0 {
		ws.userConfigPath = filepath.Join(currentAccountDir, configFileName)
	}

	// 设置配置文件路径
	config.Load(config.UserConfigPath(ws.userConfigPath), config.GlobalConfigPath(ws.globalConfigPath))

	// 加载配置
	cfg.Merge(ws.cmdConfig)
	cfg.Merge(config.GetUser())
	cfg.Merge(config.GetGlobal())
	cfg.Merge(DefaultConfig())

	if err == nil {
		cfg.Credentials = auth.Credentials{
			AccessKey: currentAccount.AccessKey,
			SecretKey: []byte(currentAccount.SecretKey),
		}
	}

	return
}

type workspace struct {
	cmdConfig        *config.Config
	workspace        string
	userConfigPath   string
	globalConfigPath string
}

func (w *workspace) initInfo() error {
	home, err := utils.GetHomePath()
	if err != nil || len(home) == 0 {
		return err
	}

	w.workspace = filepath.Join(home, workspaceName)
	w.globalConfigPath = filepath.Join(home, configFileName)
	return nil
}
