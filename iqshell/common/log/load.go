package log

import (
	"errors"
	"github.com/astaxie/beego/logs"
)

func LoadConsole(cfg Config) (err error) {
	err = progressStdoutLog.SetLogger(adapterConsole, cfg.ToJson())
	if err != nil {
		err = errors.New("load console error when set logger:" + err.Error())
		return
	}
	progressStdoutLog.SetLevel(cfg.Level)
	err = progressStdoutLog.DelLogger(logs.AdapterConsole)
	if err != nil {
		err = errors.New("load console error when del logger:" + err.Error())
	}
	return
}

func LoadFileLogger(cfg Config) (err error) {
	err = progressFileLog.SetLogger(logs.AdapterFile, cfg.ToJson())
	progressStdoutLog.DelLogger(logs.AdapterConsole)
	return
}
