package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	progressStdoutLog = new(logs.BeeLogger)
	progressFileLog   = new(logs.BeeLogger)
)

func Debug(a ...interface{}) {
	progressStdoutLog.Debug(fmt.Sprint(a...))
}

func DebugF(format string, v ...interface{}) {
	progressStdoutLog.Debug(format, v...)
}

func Info(a ...interface{}) {
	progressStdoutLog.Info(fmt.Sprint(a...))
}

func InfoF(format string, v ...interface{}) {
	progressStdoutLog.Info(format, v...)
}

func Warning(a ...interface{}) {
	progressStdoutLog.Warn(fmt.Sprint(a...))
}

func WarningF(format string, v ...interface{}) {
	progressStdoutLog.Warn(format, v...)
}

func Error(a ...interface{}) {
	progressStdoutLog.Error(fmt.Sprint(a...))
}

func ErrorF(format string, v ...interface{}) {
	progressStdoutLog.Error(format, v...)
}

func Alert(a ...interface{}) {
	progressStdoutLog.Alert(fmt.Sprint(a...))
}

func AlertF(format string, v ...interface{}) {
	progressStdoutLog.Alert(format, v...)
}
