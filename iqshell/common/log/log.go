package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var progressLog *logs.BeeLogger

func Debug(a ...interface{}) {
	DebugF(fmt.Sprint(a...))
}

func DebugF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Debug(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}
}

func Info(a ...interface{}) {
	InfoF(fmt.Sprint(a...))
}

func InfoF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Info(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}
}

func Warning(a ...interface{}) {
	WarningF(fmt.Sprint(a...))
}

func WarningF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Warn(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}
}

func Error(a ...interface{}) {
	ErrorF(fmt.Sprint(a...))
}

func ErrorF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Error(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}
}

func Alert(a ...interface{}) {
	AlertF(fmt.Sprint(a...))
}

func AlertF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Alert(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}

}
