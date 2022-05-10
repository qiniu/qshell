package log

import (
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func Prepare() *data.CodeError {
	progressLog = new(logs.BeeLogger)
	return nil
}

func LoadConsole(cfg Config) (err *data.CodeError) {
	if e := progressLog.SetLogger(adapterConsole, cfg.ToJson()); e != nil {
		return data.NewEmptyError().AppendDesc("load console error when set logger").AppendError(e)
	}

	// 日志总开关
	progressLog.SetLevel(LevelDebug)

	if e := progressLog.DelLogger(logs.AdapterConsole); e != nil {
		return data.NewEmptyError().AppendDesc("load console error when del logger").AppendError(e)
	}

	return
}

func LoadFileLogger(cfg Config) (err *data.CodeError) {
	if len(cfg.Filename) > 0 {
		if e := progressLog.SetLogger(logs.AdapterFile, cfg.ToJson()); e != nil {
			return data.NewEmptyError().AppendDesc("set file logger").AppendError(e)
		}
	}

	if !cfg.EnableStdout {
		if dErr := progressLog.DelLogger(adapterConsole); dErr != nil {
			WarningF("disable stdout error:%v", dErr)
		}
	}
	return
}
