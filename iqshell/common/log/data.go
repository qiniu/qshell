package log

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

const (
	LevelSlient  Level = -10
	LevelVerbose Level = logs.LevelDebug
	LevelDebug   Level = logs.LevelInformational
	LevelInfo    Level = logs.LevelNotice
	LevelWarning Level = logs.LevelWarning
	LevelError   Level = logs.LevelError
)

type Level int

type Config struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	Daily    bool   `json:"daily"`
	MaxDays  int    `json:"maxdays"`
}

func (c *Config) ToJson() string {
	cfgBytes, _ := json.Marshal(c)
	return string(cfgBytes)
}
