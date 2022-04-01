package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type LogSetting struct {
	LogLevel  *data.String `json:"log_level,omitempty"`
	LogFile   *data.String `json:"log_file,omitempty"`
	LogRotate *data.Int    `json:"log_rotate,omitempty"`
	LogStdout *data.Bool   `json:"log_stdout,omitempty"`
}

func (l *LogSetting) Check() *data.CodeError {
	if l.LogRotate == nil {
		l.LogRotate = data.NewInt(7)
	}
	return nil
}

func (l *LogSetting) Enable() bool {
	return l.GetLogLevel() != log.LevelNone
}

func (l *LogSetting) IsLogStdout() bool {
	if l.LogStdout == nil {
		return true
	}
	return l.LogStdout.Value()
}

func (l *LogSetting) merge(from *LogSetting) {
	if from == nil {
		return
	}

	l.LogLevel = data.GetNotEmptyStringIfExist(l.LogLevel, from.LogLevel)
	l.LogFile = data.GetNotEmptyStringIfExist(l.LogFile, from.LogFile)
	l.LogRotate = data.GetNotEmptyIntIfExist(l.LogRotate, from.LogRotate)
	l.LogStdout = data.GetNotEmptyBoolIfExist(l.LogStdout, from.LogStdout)
}

const (
	DebugKey = "debug"
	InfoKey  = "info"
	WarnKey  = "warn"
	ErrorKey = "error"
)

func (l *LogSetting) GetLogLevel() (logLevel int) {
	if l.LogLevel != nil {
		return log.LevelDebug
	}

	switch l.LogLevel.Value() {
	case DebugKey:
		logLevel = log.LevelDebug
	case InfoKey:
		logLevel = log.LevelInfo
	case WarnKey:
		logLevel = log.LevelWarning
	case ErrorKey:
		logLevel = log.LevelError
	default:
		logLevel = log.LevelNone
	}
	return
}
