package log

import (
	"github.com/astaxie/beego/logs"
)

var (
	progressStdoutLog = new(logs.BeeLogger)
	progressFileLog   = new(logs.BeeLogger)
	resultLog         = new(logs.BeeLogger)
)

func Debug(format string, v ...interface{}) {
	progressStdoutLog.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	progressStdoutLog.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	progressStdoutLog.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	progressStdoutLog.Error(format, v...)
}

func Alert(format string, v ...interface{}) {
	progressStdoutLog.Alert(format, v...)
}
