package move

import "os"

type QiniuFile struct {
	Bucket string // 空间名
	Key    string // 文件 Key
}

type LocalFile struct {
	FilePath    string   // 文件路径
	fileHandler *os.File // 文件操作句柄
}

type NetworkFile struct {
	FileUrl     string   // 文件 Url
}
