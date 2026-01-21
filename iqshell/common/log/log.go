package log

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var progressLog *logs.BeeLogger

func Debug(a ...interface{}) {
	if progressLog != nil {
		progressLog.Debug(fmt.Sprint(a...))
	} else {
		fmt.Println(a...)
	}
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
	if progressLog != nil {
		progressLog.Info(fmt.Sprint(a...))
	} else {
		fmt.Println(a...)
	}
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
	if progressLog != nil {
		progressLog.Warn(fmt.Sprint(a...))
	} else {
		fmt.Println(a...)
	}
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
	if progressLog != nil {
		progressLog.Error(fmt.Sprint(a...))
	} else {
		fmt.Println(a...)
	}
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
	if progressLog != nil {
		progressLog.Alert(fmt.Sprint(a...))
	} else {
		fmt.Println(a...)
	}
}

func AlertF(format string, v ...interface{}) {
	if progressLog != nil {
		progressLog.Alert(format, v...)
	} else {
		fmt.Printf(format, v...)
		fmt.Println("")
	}
}
