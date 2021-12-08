package iqshell

import (
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"runtime"
)

type Config struct {
	DebugEnable    bool
	DDebugEnable   bool
	ConfigFilePath string
	WorkspacePath  string
}

func Load(cfg Config) (err error) {
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
	_ = log.LoadConsole(logLevel)

	// 加载工作区
	err = workspace.Load(workspace.Workspace(cfg.WorkspacePath), workspace.UserConfigPath(cfg.ConfigFilePath))
	return
}
