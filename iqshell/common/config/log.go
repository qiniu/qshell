package config

import "github.com/astaxie/beego/logs"

type LogSetting struct {
	LogLevel   string `json:"log_level,omitempty"`
	LogFile    string `json:"log_file,omitempty"`
	LogRotate  int    `json:"log_rotate,omitempty"`
	LogStdout  bool   `json:"log_stdout,omitempty"`
}

func (l *LogSetting) merge(from *LogSetting) {
	if from == nil {
		return
	}

	if len(l.LogLevel) == 0 {
		l.LogLevel = from.LogLevel
	}

	if len(l.LogFile) == 0 {
		l.LogFile = from.LogFile
	}

	if l.LogRotate == 0 {
		l.LogRotate = from.LogRotate
	}

	if !l.LogStdout {
		l.LogStdout = from.LogStdout
	}
}

func (l *LogSetting)GetLogLevel() (logLevel int) {
	switch l.LogLevel {
	case "debug":
		logLevel = logs.LevelDebug
	case "info":
		logLevel = logs.LevelInfo
	case "warn":
		logLevel = logs.LevelWarning
	case "error":
		logLevel = logs.LevelError
	default:
		logLevel = logs.LevelInfo
	}
	return
}