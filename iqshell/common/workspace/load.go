package workspace

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"path/filepath"

	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type LoadInfo struct {
	UserConfigPath   string
	CmdConfig        *config.Config
	WorkspacePath    string
	JobPathBuilder   func(cmdPath string) string
	globalConfigPath string
}

// Load 加载工作环境
func Load(info LoadInfo) (err *data.CodeError) {
	err = info.initInfo()
	if err != nil {
		return
	}

	err = config.LoadGlobalConfig(info.globalConfigPath)
	if err != nil {
		log.ErrorF("load config error:%v", err)
		return
	}

	// 检查工作目录
	if len(info.WorkspacePath) == 0 {
		err = data.NewEmptyError().AppendDesc("can't get home dir")
		return
	}
	workspaceDir = info.WorkspacePath
	log.Debug("workspace:" + workspaceDir)

	err = utils.CreateDirIfNotExist(workspaceDir)
	if err != nil {
		log.ErrorF("create workspace dir error:%v", err)
		return
	}

	// 加载账户
	accountDBPath := filepath.Join(workspaceDir, usersDBName)
	accountPath := filepath.Join(workspaceDir, currentUserFileName)
	oldAccountPath := filepath.Join(workspaceDir, oldUserFileName)
	err = account.Load(account.LoadInfo{
		AccountPath:    accountPath,
		OldAccountPath: oldAccountPath,
		AccountDBPath:  accountDBPath,
	})
	if err != nil {
		log.ErrorF("load account error:%v", err)
		return
	}

	if len(info.UserConfigPath) > 0 {
		// 用户配置了路径，使用用户的路径加载配置
		err = config.LoadUserConfig(info.UserConfigPath)
		if err != nil {
			log.ErrorF("load config error:%v", err)
			return
		}

		loadUserInfo()
	} else {
		loadUserInfo()
		info.UserConfigPath = filepath.Join(userDir, configFileName)

		err = config.LoadUserConfig(info.UserConfigPath)
		if err != nil {
			log.ErrorF("load config error:%v", err)
			return
		}
	}

	// 加载配置
	resultCfg := &config.Config{}
	resultCfg.Merge(info.CmdConfig)
	resultCfg.Merge(config.GetUser())
	resultCfg.Merge(config.GetGlobal())
	resultCfg.Merge(defaultConfig())
	cfg = resultCfg

	log.DebugF("cmd    config:\n%v", info.CmdConfig)
	log.DebugF("user   config(%s):\n%v", info.UserConfigPath, config.GetUser())
	log.DebugF("global config(%s):\n%v", info.globalConfigPath, config.GetGlobal())
	log.DebugF("final  config:\n%v", cfg)

	err = checkConfig(cfg)
	if err != nil {
		return
	}

	// 配置 Job path
	jobDir = filepath.Join(userDir, info.CmdConfig.CmdId)
	if info.JobPathBuilder != nil {
		jobDir = info.JobPathBuilder(jobDir)
	}
	err = utils.CreateDirIfNotExist(jobDir)
	if err != nil {
		return data.NewEmptyError().AppendDescF("create job dir error:%v", err)
	}
	log.DebugF("job dir:%s", jobDir)

	// uc 缓存路径, 实际路径在用户目录下，不存在 uc_cache
	storage.SetRegionCachePath(filepath.Join(userDir, "uc_cache"))

	// 在工作区加载之后监听
	observerCmdInterrupt()

	return
}

func (w *LoadInfo) initInfo() *data.CodeError {
	home, err := utils.GetHomePath()
	if err != nil {
		return data.NewEmptyError().AppendDescF("get home path error:%v", err)
	}
	if len(w.WorkspacePath) == 0 {
		w.WorkspacePath = filepath.Join(home, workspaceName)
	}
	// 全局配置文件路径，兼容老版本，位置在用户目录下
	w.globalConfigPath = filepath.Join(home, configFileName)
	return nil
}

func loadUserInfo() {
	acc, err := account.GetAccount()
	if err == nil {
		currentAccount = &acc
		accountName := acc.Name
		if len(accountName) == 0 {
			accountName = currentAccount.AccessKey
		}
		log.DebugF("current user name:%s", accountName)

		userDir = filepath.Join(workspaceDir, usersDirName, accountName)

		// 配置 config 的 Credentials
		cfg.Credentials = &auth.Credentials{
			AccessKey: acc.AccessKey,
			SecretKey: []byte(acc.SecretKey),
		}
	} else {
		userDir = filepath.Join(workspaceDir, usersDirName, defaultUserDirName)
	}

	log.DebugF("user dir:%s", userDir)
}
