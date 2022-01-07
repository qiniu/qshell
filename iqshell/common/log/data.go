package log

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

const (
	LevelAlert   Level = logs.LevelAlert
	LevelError   Level = logs.LevelError
	LevelWarning Level = logs.LevelWarning
	LevelInfo    Level = logs.LevelInformational
	LevelDebug   Level = logs.LevelDebug
)

type Level int

type Config struct {
	Filename       string `json:"filename"`
	Level          int    `json:"level"`
	Daily          bool   `json:"daily"`
	StdOutColorful bool   `json:"color"`
	MaxDays        int    `json:"maxdays"`
}

func (c *Config) ToJson() string {
	cfgBytes, _ := json.Marshal(c)
	return string(cfgBytes)
}
