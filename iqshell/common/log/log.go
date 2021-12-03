package log

import (
	"encoding/json"
)

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
