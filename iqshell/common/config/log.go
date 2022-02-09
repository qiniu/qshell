package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type LogSetting struct {
	LogLevel  string `json:"log_level,omitempty"`
	LogFile   string `json:"log_file,omitempty"`
	LogRotate int    `json:"log_rotate,omitempty"`
	LogStdout string `json:"log_stdout,omitempty"`
}

func (l *LogSetting) IsLogStdout() bool {
	return l.LogStdout != data.FalseString
}

func (l *LogSetting) merge(from *LogSetting) {
	if from == nil {
		return
	}

	l.LogLevel = utils.GetNotEmptyStringIfExist(l.LogLevel, from.LogLevel)
	l.LogFile = utils.GetNotEmptyStringIfExist(l.LogFile, from.LogFile)
	l.LogRotate = utils.GetNotZeroIntIfExist(l.LogRotate, from.LogRotate)
	l.LogStdout = utils.GetNotEmptyStringIfExist(l.LogStdout, from.LogStdout)
}

const (
	DebugKey = "debug"
	InfoKey  = "info"
	WarnKey  = "warn"
	ErrorKey = "error"
)

func (l *LogSetting) GetLogLevel() (logLevel int) {
	switch l.LogLevel {
	case DebugKey:
		logLevel = log.LevelDebug
	case InfoKey:
		logLevel = log.LevelInfo
	case WarnKey:
		logLevel = log.LevelWarning
	case ErrorKey:
		logLevel = log.LevelError
	default:
		logLevel = log.LevelDebug
	}
	return
}
