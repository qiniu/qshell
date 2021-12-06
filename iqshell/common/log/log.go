package log

import (
	"github.com/astaxie/beego/logs"
)

var (
	progressLog = new(logs.BeeLogger)
	resultLog   = new(logs.BeeLogger)
)

func Debug(format string, v ...interface{}) {
	progressLog.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	progressLog.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	progressLog.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	progressLog.Error(format, v...)
}
