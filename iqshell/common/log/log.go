package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	progressLog = new(logs.BeeLogger)
)

func Debug(a ...interface{}) {
	progressLog.Debug(fmt.Sprint(a...))
}

func DebugF(format string, v ...interface{}) {
	progressLog.Debug(format, v...)
}

func Info(a ...interface{}) {
	progressLog.Info(fmt.Sprint(a...))
}

func InfoF(format string, v ...interface{}) {
	progressLog.Info(format, v...)
}

func Warning(a ...interface{}) {
	progressLog.Warn(fmt.Sprint(a...))
}

func WarningF(format string, v ...interface{}) {
	progressLog.Warn(format, v...)
}

func Error(a ...interface{}) {
	progressLog.Error(fmt.Sprint(a...))
}

func ErrorF(format string, v ...interface{}) {
	progressLog.Error(format, v...)
}

func Alert(a ...interface{}) {
	progressLog.Alert(fmt.Sprint(a...))
}

func AlertF(format string, v ...interface{}) {
	progressLog.Alert(format, v...)
}
