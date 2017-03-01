package qshell

import (
	"encoding/json"
)

type BeeLogConfig struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	Daily    bool   `json:"daily"`
	MaxDays  int    `json:"maxdays"`
}

func (c *BeeLogConfig) ToJson() string {
	cfgBytes, _ := json.Marshal(c)
	return string(cfgBytes)
}
