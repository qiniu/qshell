package download

import (
	"errors"
)

type ApiInfo struct {
	Url                 string // 文件下载的 url 【必填】
	Domain              string // 文件下载的 domain 【必填】
	ToFile              string // 文件保存的路径 【必填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件 【必填】
	Referer             string // 请求 header 中的 Referer 【选填】
	FileEncoding        string // 文件编码方式 【选填】
	Bucket              string // 文件所在 bucket，用于验证 hash 【选填】
	Key                 string // 文件被保存的 key，用于验证 hash 【选填】
	FileHash            string // 文件 hash，有值则会检测 hash 【选填】
	FileSize            int64  // 文件大小，有值则会检测文件大小 【选填】
	FileModifyTime      int64  // 文件修改时间 【选填】
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info ApiInfo) (file string, err error) {
	if len(info.ToFile) == 0 {
		return file, errors.New("the filename saved after downloading is empty")
	}

	downloader := &Downloader{
		Url:          info.Url,
		ToFile:       info.ToFile,
		Domain:       info.Domain,
		Referer:      info.Referer,
		FileEncoding: info.FileEncoding,
	}
	err = downloader.Download()
	if err != nil {
		return file, err
	}

	err = (&Checker{
		File:     downloader.ToFile,
		Bucket:   info.Bucket,
		Key:      info.Key,
		FileHash: info.FileHash,
		FileSize: info.FileSize,
	}).Check()

	return info.ToFile, err
}
