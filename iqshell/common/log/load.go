package log

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func Prepare() error {
	progressLog = new(logs.BeeLogger)
	return nil
}

func LoadConsole(cfg Config) (err error) {
	err = progressLog.SetLogger(adapterConsole, cfg.ToJson())
	if err != nil {
		err = errors.New("load console error when set logger:" + err.Error())
		return
	}
	// 日志总开关
	progressLog.SetLevel(LevelDebug)
	err = progressLog.DelLogger(logs.AdapterConsole)
	if err != nil {
		err = errors.New("load console error when del logger:" + err.Error())
	}
	return
}

func LoadFileLogger(cfg Config) (err error) {
	if len(cfg.Filename) > 0 {
		err = progressLog.SetLogger(logs.AdapterFile, cfg.ToJson())
		if err != nil {
			err = fmt.Errorf("set file logger error:%v", err)
		}
	}

	if !cfg.EnableStdout {
		if dErr := progressLog.DelLogger(adapterConsole); dErr != nil {
			WarningF("disable stdout error:%v", dErr)
		}
	}
	return
}
