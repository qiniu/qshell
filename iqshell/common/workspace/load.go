package workspace

import (
	"errors"
	"path/filepath"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type LoadInfo struct {
	UserConfigPath   string
	CmdConfig        *config.Config
	WorkspacePath    string
	globalConfigPath string
}

// Load 加载工作环境
func Load(info LoadInfo) (err error) {
	err = info.initInfo()
	if err != nil {
		return
	}

	// 检查工作目录
	if len(info.WorkspacePath) == 0 {
		err = errors.New("can't get home dir")
		return
	}
	workspacePath = info.WorkspacePath
	log.Debug("workspace:" + workspacePath)

	err = utils.CreateDirIfNotExist(workspacePath)
	if err != nil {
		return
	}

	// 加载账户
	accountDBPath := filepath.Join(workspacePath, usersDBName)
	accountPath := filepath.Join(workspacePath, currentUserFileName)
	oldAccountPath := filepath.Join(workspacePath, oldUserFileName)
	err = account.Load(account.LoadInfo{
		AccountPath:    accountPath,
		OldAccountPath: oldAccountPath,
		AccountDBPath:  accountDBPath,
	})
	if err != nil {
		log.ErrorF("load account error:%v", err)
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

		currentAccountDir = filepath.Join(workspacePath, accountName)
		err := utils.CreateDirIfNotExist(currentAccountDir)
		if err != nil {
			return errors.New("create user dir error:" + err.Error())
		}
	}

	// 检查用户配置，用户配置可能被指定，如果未指定则使用用户目录下配置
	if len(info.UserConfigPath) == 0 {
		info.UserConfigPath = filepath.Join(currentAccountDir, configFileName)
	}

	// 设置配置文件路径
	err = config.Load(config.LoadInfo{
		UserConfigPath:   info.UserConfigPath,
		GlobalConfigPath: info.globalConfigPath,
	})
	if err != nil {
		log.ErrorF("load config error:%v", err)
		return
	}

	// 加载配置
	resultCfg := &config.Config{}
	resultCfg.Merge(info.CmdConfig)
	resultCfg.Merge(config.GetUser())
	resultCfg.Merge(config.GetGlobal())
	resultCfg.Merge(defaultConfig())
	cfg = resultCfg

	log.DebugF("cmd    config:\n%v", info.CmdConfig)
	log.DebugF("user   config:\n%v", config.GetUser())
	log.DebugF("global config:\n%v", config.GetGlobal())
	log.DebugF("final  config:\n%v", cfg)

	err = checkConfig(cfg)
	if err != nil {
		return
	}

	cfg.Credentials = &auth.Credentials{
		AccessKey: currentAccount.AccessKey,
		SecretKey: []byte(currentAccount.SecretKey),
	}

	// 在工作区加载之后监听
	observerCmdInterrupt()

	return
}

func (w *LoadInfo) initInfo() error {
	home, err := utils.GetHomePath()
	if err != nil {
		return errors.New("get home path error:" + err.Error())
	}
	if len(w.WorkspacePath) == 0 {
		w.WorkspacePath = filepath.Join(home, workspaceName)
	}
	// 全局配置文件路径，兼容老版本，位置在用户目录下
	w.globalConfigPath = filepath.Join(home, configFileName)
	return nil
}
