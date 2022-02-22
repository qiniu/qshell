package iqshell

import (
	"errors"
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

type Config struct {
	DebugEnable        bool   // 开启命令行的调试模式
	DDebugEnable       bool   // go SDK client 和命令行开启调试模式
	ConfigFilePath     string // 配置文件路径，用户可以指定配置文件
	Local              bool   // 是否使用当前文件夹作为工作区
	StdoutColorful     bool   // 控制台输出是否多彩
	UploadConfigFile   string // 上传配置文件
	DownloadConfigFile string // 下载配置文件
	CmdCfg             config.Config
}

func Load(cfg Config) error {

	// 加载 log
	logLevel := log.LevelInfo
	if cfg.DebugEnable || cfg.DDebugEnable {
		logLevel = log.LevelDebug
	}
	if cfg.DDebugEnable {
		client.TurnOnDebug()
	}

	// 加载本地输出
	_ = log.LoadConsole(log.Config{
		Level:          logLevel,
		StdOutColorful: cfg.StdoutColorful,
	})

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

	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 配置 user agent
	storage.UserAgent = utils.UserAgent()

	// 合并上传配置
	if len(cfg.UploadConfigFile) > 0 {
		if err := utils.UnMarshalFromFile(cfg.UploadConfigFile, &cfg.CmdCfg.Up); err != nil {
			return fmt.Errorf("read upload config error:%v config file:%s", err, cfg.UploadConfigFile)
		}
	}

	// 合并下载配置
	if len(cfg.DownloadConfigFile) > 0 {
		if err := utils.UnMarshalFromFile(cfg.DownloadConfigFile, &cfg.CmdCfg.Download); err != nil {
			return fmt.Errorf("read download config error:%v config file:%s", err, cfg.UploadConfigFile)
		}
	}

	// 加载工作区
	if err := workspace.Load(workspace.LoadInfo{
		CmdConfig:      &cfg.CmdCfg,
		WorkspacePath:  workspacePath,
		UserConfigPath: cfg.ConfigFilePath,
	}); err != nil {
		return err
	}

	if ls := workspace.GetLogConfig(); ls != nil {
		err := utils.CreateFileDirIfNotExist(ls.LogFile)
		if err != nil {
			return errors.New("create log file error:" + err.Error())
		}
		_ = log.LoadFileLogger(log.Config{
			Filename:       ls.LogFile,
			Level:          ls.GetLogLevel(),
			Daily:          true,
			StdOutColorful: false,
			EnableStdout:   ls.IsLogStdout(),
			MaxDays:        ls.LogRotate,
		})
	}

	return nil
}
