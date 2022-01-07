package iqshell

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"runtime"
)

type Runnable interface {
	Run()
}

type Config struct {
	DebugEnable    bool   // 开启命令行的调试模式
	DDebugEnable   bool   // go SDK client 和命令行开启调试模式
	ConfigFilePath string // 配置文件路径，用户可以指定配置文件
	Local          bool   // 是否使用当前文件夹作为工作区
	CmdCfg         config.Config
}

func Load(cfg Config) error {

	workspacePath := ""
	if cfg.Local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		workspacePath = dir
	}

	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())

	storage.UserAgent = utils.UserAgent()

	// 加载 log
	logLevel := log.LevelInfo
	if cfg.DebugEnable || cfg.DDebugEnable {
		logLevel = log.LevelDebug
	}
	if cfg.DDebugEnable {
		client.TurnOnDebug()
	}

	//TODO: 此处逻辑处理
	if cfg.CmdCfg.Log.Level == 0 {
		cfg.CmdCfg.Log.Level = int(logLevel)
	}
	cfg.CmdCfg.Log.StdOutColorful = false
	_ = log.LoadConsole(cfg.CmdCfg.Log)

	// 加载工作区
	if err := workspace.Load(workspace.Workspace(workspacePath),
		workspace.UserConfigPath(cfg.ConfigFilePath),
		workspace.CmdConfig(&cfg.CmdCfg)); err != nil {
		return err
	}

	return nil
}
