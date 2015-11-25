package qshell

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Hex(from string) string {
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(from))
	return hex.EncodeToString(md5Hasher.Sum(nil))
}
