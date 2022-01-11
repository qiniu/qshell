package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"strconv"
)

type BatchDownloadInfo struct {
	Bucket              string // 下载的 bucket 【必填】
	Domain              string // 文件下载的 domain 【必填】
	ToFile              string // 文件保存的路径 【必填】
	IsPublic            bool   // 是否是公有云 【必填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件 【必填】
	Referer             string // 请求 header 中的 Referer 【选填】
	FileEncoding        string // 文件编码方式 【选填】
	GroupInfo           group.Info
}

func BatchDownload(info BatchDownloadInfo) {
	handler, err := group.NewHandler(info.GroupInfo)
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		_ = handler.Scanner().Close()
	}()

	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, false
		}

		items := utils.SplitString(line, info.GroupInfo.ItemSeparate)
		if len(items) < 4 {
			log.ErrorF("invalid line, line must contain key fileSize fileHash and fileModifyTime:%s", line)
			handler.Export().Fail().ExportF("%s: error:%s", line, "line must contain key fileSize fileHash and fileModifyTime")
			return nil, true
		}

		fileKey := items[0]
		fileSize, err := strconv.ParseInt(items[1], 10, 64)
		if err != nil {
			log.ErrorF("invalid line, get file size error:%s", line)
			handler.Export().Fail().ExportF("%s: get file size error:%s", line, err)
			return nil, true
		}

		fileHash := items[2]
		fileModifyTime, err := strconv.ParseInt(items[3], 10, 64)
		if err != nil {
			log.ErrorF("invalid line, get file modify time error:%s", line)
			handler.Export().Fail().ExportF("%s: get file modify time error:%s", line, err)
			return nil, true
		}

		return DownloadInfo{
			ApiInfo: download.ApiInfo{
				Url:                 "",
				Domain:              info.Domain,
				ToFile:              info.ToFile,
				RemoveFileWhenError: true,
				Referer:             info.Referer,
				FileEncoding:        info.FileEncoding,
				Bucket:              info.Bucket,
				Key:                 fileKey,
				FileHash:            fileHash,
				FileSize:            fileSize,
				FileModifyTime:      fileModifyTime,
			},
			IsPublic: info.IsPublic,
		}, true
	}).DoWork(func(work work.Work) (work.Result, error) {
		info := work.(DownloadInfo)
		file, err := downloadFile(info)
		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(download.ApiInfo)
		file := result.(string)
		handler.Export().Success().ExportF("download success, [%s:%s] => %s", apiInfo.Bucket, apiInfo.Key, file)
	}).OnWorkError(func(work work.Work, err error) {
		apiInfo := work.(download.ApiInfo)
		handler.Export().Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* key fileSize fileHash and fileModifyTime */
			apiInfo.Key, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileHash, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()
}
