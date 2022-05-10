package utils

import (
	"runtime"
	"strings"
)

func IsWindowsOS() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func IsGBKEncoding(encoding string) bool {
	return strings.ToLower(encoding) == "gbk"
}
