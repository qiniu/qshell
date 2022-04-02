package utils

import (
	"math/rand"
	"strings"
)

const (
	// ASCII英文字母
	alphaList = "abcdefghijklmnopqrstuvwxyz"
)

// CreateRandString 生成随机的字符串
func CreateRandString(num int) (rcode string) {
	if num <= 0 || num > len(alphaList) {
		rcode = ""
		return
	}

	buffer := make([]byte, num)
	_, err := rand.Read(buffer)
	if err != nil {
		rcode = ""
		return
	}

	for _, b := range buffer {
		index := int(b) / len(alphaList)
		rcode += string(alphaList[index])
	}

	return
}

func SplitString(line, sep string) []string {
	if len(sep) == 0 {
		//strings.TrimSpace(sep) == ""
		return strings.Fields(line)
	}
	return strings.Split(line, sep)
}
