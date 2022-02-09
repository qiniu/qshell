package log

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
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
		err = progressStdoutLog.SetLogger(logs.AdapterFile, cfg.ToJson())
		if err != nil {
			err = fmt.Errorf("set file logger error:%v", err)
		}
	}

	if !cfg.EnableStdout {
		if dErr := progressStdoutLog.DelLogger(adapterConsole); dErr != nil {
			log.Warnf("disable stdout error:%v", dErr)
		}
	}
	return
}
