package iqshell

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"runtime"
)

type Config struct {
	Document       bool   // 是否展示 document
	DebugEnable    bool   // 开启命令行的调试模式
	DDebugEnable   bool   // go SDK client 和命令行开启调试模式
	ConfigFilePath string // 配置文件路径，用户可以指定配置文件
	Local          bool   // 是否使用当前文件夹作为工作区
	StdoutColorful bool   // 控制台输出是否多彩
	CmdCfg         config.Config
}

type CheckAndLoadInfo struct {
	Checker           data.Checker
	BeforeLoadFileLog func()
	AfterLoadFileLog  func()
}

func CheckAndLoad(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if !Load(cfg, info) {
		return false
	}
	return Check(cfg, info)
}

func Load(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if ShowDocumentIfNeeded(cfg) {
		return false
	}
	if !LoadBase(cfg) {
		return false
	}
	if !LoadWorkspace(cfg) {
		return false
	}
	if info.BeforeLoadFileLog != nil {
		info.BeforeLoadFileLog()
	}
	shouldContinue = LoadFileLog(cfg)
	if info.AfterLoadFileLog != nil {
		info.AfterLoadFileLog()
	}
	if !shouldContinue {
		return false
	}

	return true
}

func Check(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if info.Checker != nil {
		if err := info.Checker.Check(); err != nil {
			log.ErrorF("check error: %v", err)
			return false
		}
	}
	return true
}

func ShowDocumentIfNeeded(cfg *Config) bool {
	if !cfg.Document {
		return false
	}
	docs.ShowCmdDocument(cfg.CmdCfg.CmdId)
	return true
}

func LoadBase(cfg *Config) (shouldContinue bool) {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 配置 user agent
	storage.UserAgent = utils.UserAgent()

	// 加载 log
	logLevel := log.LevelInfo
	if cfg.DebugEnable || cfg.DDebugEnable {
		logLevel = log.LevelDebug
	}
	if cfg.DDebugEnable {
		client.TurnOnDebug()
	}

	// 加载本地输出
	_ = log.Prepare()
	_ = log.LoadConsole(log.Config{
		Level:          logLevel,
		StdOutColorful: cfg.StdoutColorful,
	})
	return true
}

func LoadWorkspace(cfg *Config) (shouldContinue bool) {
	// 获取工作目录
	workspacePath := ""
	if cfg.Local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		workspacePath = dir
	}

	// 加载工作区
	if err := workspace.Load(workspace.LoadInfo{
		CmdConfig:      &cfg.CmdCfg,
		WorkspacePath:  workspacePath,
		UserConfigPath: cfg.ConfigFilePath,
	}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load workspace error: %v\n", err)
		return false
	}
	return true
}

func LoadFileLog(cfg *Config) (shouldContinue bool) {
	// 配置日志文件输出
	if ls := workspace.GetLogConfig(); ls != nil && ls.Enable() && data.NotEmpty(ls.LogFile) {
		err := utils.CreateFileDirIfNotExist(ls.LogFile.Value())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "load file log, create log file error: %v\n", err)
			return false
		}
		_ = log.LoadFileLogger(log.Config{
			Filename:       ls.LogFile.Value(),
			Level:          ls.GetLogLevel(),
			Daily:          true,
			StdOutColorful: false,
			EnableStdout:   ls.IsLogStdout(),
			MaxDays:        ls.LogRotate.Value(),
		})
		log.AlertF("Writing log to file:%s \n\n", workspace.GetConfig().Log.LogFile.Value())
	} else {
		log.DebugF("log file not enable, log level:%s \n\n", workspace.GetConfig().Log.LogLevel.Value())
	}
	return true
}
