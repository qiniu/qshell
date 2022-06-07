package local

import "os"

type localFile struct {
	filePath    string   // 文件路径
	fileHandler *os.File // 文件操作句柄
}


