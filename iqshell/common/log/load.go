package log

import (
	"errors"
	"fmt"
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
	if len(cfg.Filename) > 0 {
		err = progressFileLog.SetLogger(logs.AdapterFile, cfg.ToJson())
		if err != nil {
			return fmt.Errorf("progressFileLog set logger error:%v", err)
		}
	}

	_ = progressStdoutLog.DelLogger(logs.AdapterConsole)
	return
}
