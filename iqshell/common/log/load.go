package log

import "github.com/astaxie/beego/logs"

func LoadConsole(logLevel Level) error {
	progressLog.SetLogger(logs.AdapterConsole)
	progressLog.SetLevel(int(logLevel))
	// resultLog.SetLogger(logs.AdapterFile, log.Config{
	// 	Filename: downConfig.LogFile,
	// 	Level:    logLevel,
	// 	Daily:    true,
	// 	MaxDays:  logRotate,
	// })
	return nil
}

func LoadFileLogger(cfg Config) (err error) {
	err = progressLog.SetLogger(logs.AdapterFile, cfg.ToJson())
	return
}
