package iqshell

import (
	"encoding/json"
)

type BeeLogConfig struct {
	Filename     string `json:"filename"`
	Level        int    `json:"level"`
	Daily        bool   `json:"daily"`
	MaxDays      int    `json:"maxdays"`
	MinFileCount int    `json:"min_file_count"` // 最少文件数
}

func (c *BeeLogConfig) ToJson() string {
	cfgBytes, _ := json.Marshal(c)
	return string(cfgBytes)
}
