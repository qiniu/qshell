package log

import "github.com/astaxie/beego/logs"

func LoadConsole(logLevel Level) error {
	progressStdoutLog.SetLogger(adapterConsole)
	progressStdoutLog.SetLevel(int(logLevel))
	progressStdoutLog.DelLogger(logs.AdapterConsole)
	// resultLog.SetLogger(logs.AdapterFile, log.Config{
	// 	Filename: downConfig.LogFile,
	// 	Level:    logLevel,
	// 	Daily:    true,
	// 	MaxDays:  logRotate,
	// })
	return nil
}

func LoadFileLogger(cfg Config) (err error) {
	err = progressFileLog.SetLogger(logs.AdapterFile, cfg.ToJson())
	progressStdoutLog.DelLogger(logs.AdapterConsole)
	return
}
