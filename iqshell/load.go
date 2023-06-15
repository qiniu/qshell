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
	"github.com/qiniu/qshell/v2/iqshell/common/version"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Document       bool                        // 是否展示 document
	Silence        bool                        // 开启命令行的静默模式，只输出 Error 和 Warning
	DebugEnable    bool                        // 开启命令行的调试模式
	DDebugEnable   bool                        // go SDK client 和命令行开启调试模式
	ConfigFilePath string                      // 配置文件路径，用户可以指定配置文件
	Local          bool                        // 是否使用当前文件夹作为工作区
	StdoutColorful bool                        // 控制台输出是否多彩
	JobPathBuilder func(cmdPath string) string // job 路径生成器
	CmdCfg         config.Config
}

type CheckAndLoadInfo struct {
	Checker           data.Checker
	BeforeLoadFileLog func()
	AfterLoadFileLog  func()
}

func CheckAndLoad(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if ShowDocumentIfNeeded(cfg) {
		return false
	}
	if !load(cfg, info) {
		return false
	}
	return Check(cfg, info)
}

func load(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if !loadBase(cfg) {
		data.SetCmdStatusError()
		return false
	}

	if !loadWorkspace(cfg) {
		data.SetCmdStatusError()
		return false
	}

	if info.BeforeLoadFileLog != nil {
		info.BeforeLoadFileLog()
	}
	shouldContinue = loadFileLog(cfg)
	if info.AfterLoadFileLog != nil {
		info.AfterLoadFileLog()
	}
	if !shouldContinue {
		data.SetCmdStatusError()
		return false
	}

	outputSomeInformationForDebug()
	return true
}

func Check(cfg *Config, info CheckAndLoadInfo) (shouldContinue bool) {
	if info.Checker != nil {
		if err := info.Checker.Check(); err != nil {
			log.ErrorF("check error: %v", err)
			data.SetCmdStatusError()
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

func loadBase(cfg *Config) (shouldContinue bool) {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 配置 user agent
	storage.UserAgent = utils.UserAgent()

	// 加载 log
	logLevel := log.LevelInfo
	if cfg.DDebugEnable {
		logLevel = log.LevelDebug
		// 深度 Debug ，client 开启日志模式
		client.TurnOnDebug()
	} else if cfg.DebugEnable {
		logLevel = log.LevelDebug
	} else if cfg.Silence {
		logLevel = log.LevelWarning
	}

	// 加载本地输出
	_ = log.Prepare()
	_ = log.LoadConsole(log.Config{
		Level:          logLevel,
		StdOutColorful: cfg.StdoutColorful,
	})
	return true
}

func loadWorkspace(cfg *Config) (shouldContinue bool) {
	// 获取工作目录
	workspacePath := ""
	if cfg.Local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(data.StatusError)
		}
		workspacePath = dir
	}

	// 加载工作区
	if err := workspace.Load(workspace.LoadInfo{
		CmdConfig:      &cfg.CmdCfg,
		WorkspacePath:  workspacePath,
		UserConfigPath: cfg.ConfigFilePath,
		JobPathBuilder: cfg.JobPathBuilder,
	}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load workspace error: %v\n", err)
		return false
	}
	return true
}

func loadFileLog(cfg *Config) (shouldContinue bool) {
	// 配置日志文件输出
	if ls := workspace.GetLogConfig(); ls != nil && ls.Enable() {
		if data.Empty(ls.LogFile) {
			ls.LogFile = data.NewString(filepath.Join(workspace.GetJobDir(), "log.txt"))
		}
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
		log.AlertF("Writing log to file:%s", ls.LogFile.Value())
	} else {
		log.DebugF("log file not enable, log level:%s", workspace.GetConfig().Log.LogLevel.Value())
	}
	return true
}

func outputSomeInformationForDebug() {
	log.DebugF("%-15s:%s", "Version", version.Version())
	log.DebugF("%-15s:%s", "UserName", workspace.GetUserName())
	log.DebugF("%-15s:%s", "Workspace", workspace.GetWorkspace())
	log.DebugF("%-15s:%s", "UserDir", workspace.GetUserDir())
	log.DebugF("%-15s:%s", "JobDir", workspace.GetJobDir())
	log.DebugF("%-15s:%s", "OS", runtime.GOOS)
	log.DebugF("%-15s:%s", "OSArch", runtime.GOARCH)
	log.DebugF("%-15s:%d", "NumCpu", runtime.NumCPU())
	log.Debug("")
}
