package download

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type ApiInfo struct {
	Url            string // 文件下载的 url 【必填】
	Domain         string // 文件下载的 domain 【必填】
	ToFile         string // 文件保存的路径 【必填】
	StatusDBPath   string // 下载状态缓存的 db 路径 【必填】
	Referer        string // 请求 header 中的 Referer 【选填】
	FileEncoding   string // 文件编码方式 【选填】
	Bucket         string // 文件所在 bucket，用于验证 hash 【选填】
	Key            string // 文件被保存的 key，用于验证 hash 【选填】
	FileHash       string // 文件 hash，有值则会检测 hash 【选填】
	FileSize       int64  // 文件大小，有值则会检测文件大小 【选填】
	FileModifyTime int64  // 文件修改时间 【选填】
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info ApiInfo) (file string, err error) {
	if len(info.ToFile) == 0 {
		return file, errors.New("the filename saved after downloading is empty")
	}

	f, err := createDownloadFiles(info.ToFile, info.FileEncoding)
	if err != nil {
		return file, err
	}

	// 检查文件是否已存在，如果存在是否符合预期
	dbChecker := &dbHandler{
		DBFilePath:           info.StatusDBPath,
		FilePath:             f.toAbsFile,
		FileHash:             info.FileHash,
		FileSize:             info.FileSize,
		FileServerModifyTime: info.FileModifyTime,
	}
	err = dbChecker.init()
	if err != nil {
		return file, err
	}

	exist, err := dbChecker.checkInfoOfDB()
	if err != nil {
		return file, err
	}
	if exist {
		log.WarningF("Local file:`%s` exists for key:`%s`", f.toAbsFile, info.Key)
		return f.toAbsFile, nil
	}

	// 下载
	downloader := &Downloader{
		files:   f,
		Url:     info.Url,
		Domain:  info.Domain,
		Referer: info.Referer,
	}
	err = downloader.Download()
	if err != nil {
		return file, err
	}

	// 下载后检查数据
	err = (&LocalFileInfo{
		File:                f.toAbsFile,
		Bucket:              info.Bucket,
		Key:                 info.Key,
		FileHash:            info.FileHash,
		FileSize:            info.FileSize,
		RemoveFileWhenError: true,
	}).CheckDownloadFile()
	if err != nil {
		return file, err
	}

	return info.ToFile, err
}
