package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
	"strconv"
)

type BatchDownloadInfo struct {
	GroupInfo group.Info
}

func BatchDownload(info BatchDownloadInfo) {
	downloadCfg := workspace.GetConfig().Download
	downloadDomain := downloadCfg.DownloadDomain()
	if len(downloadDomain) == 0 {
		log.Error("download domain can't be empty, you can set cdn_domain or io_host")
		return
	}

	dbDir := ""
	if len(downloadCfg.RecordRoot) == 0 {
		dbDir = filepath.Join(workspace.GetWorkspace(), "download")
	} else {
		dbDir = filepath.Join(downloadCfg.RecordRoot, "download")
	}
	log.InfoF("download db dir:%s", dbDir)

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
				Url:            "", // downloadFile 时会自动创建
				Domain:         downloadDomain,
				ToFile:         filepath.Join(downloadCfg.DestDir, fileKey),
				StatusDBPath:   dbDir,
				Referer:        downloadCfg.Referer,
				FileEncoding:   downloadCfg.FileEncoding,
				Bucket:         downloadCfg.Bucket,
				Key:            fileKey,
				FileHash:       fileHash,
				FileSize:       fileSize,
				FileModifyTime: fileModifyTime,
			},
			IsPublic: downloadCfg.Public,
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
