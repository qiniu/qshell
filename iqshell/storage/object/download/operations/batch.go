package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
)

type BatchDownloadInfo struct {
	GroupInfo group.Info
}

func BatchDownload(info BatchDownloadInfo) {
	downloadCfg := workspace.GetConfig().Download
	info.GroupInfo.InputFile = downloadCfg.KeyFile
	info.GroupInfo.ItemSeparate = "\t"

	downloadDomain := downloadCfg.DownloadDomain()
	if len(downloadDomain) == 0 {
		downloadDomain, _ = bucket.DomainOfBucket(downloadCfg.Bucket)
	}
	if len(downloadDomain) == 0 {
		log.Error("download domain can't be empty, you can set cdn_domain or io_host")
		return
	}

	jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", downloadCfg.DestDir, downloadCfg.Bucket, downloadCfg.KeyFile))
	dbPath := ""
	if len(downloadCfg.RecordRoot) == 0 {
		dbPath = filepath.Join(workspace.GetWorkspace(), "download", jobId, ".list")
	} else {
		dbPath = filepath.Join(downloadCfg.RecordRoot, "download", jobId, ".list")
	}
	log.InfoF("download db dir:%s", dbPath)

	export, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.GroupInfo.SuccessExportFilePath,
		FailExportFilePath:     info.GroupInfo.FailExportFilePath,
		OverrideExportFilePath: info.GroupInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.Error(err)
		return
	}

	ds, err := newDownloadScanner(downloadCfg.KeyFile, info.GroupInfo.ItemSeparate, downloadCfg.Bucket, export)
	if err != nil {
		log.Error(err)
		return
	}

	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		return ds.scan()
	}).DoWork(func(work work.Work) (work.Result, error) {
		apiInfo := work.(*download.ApiInfo)
		apiInfo.Url = "" // downloadFile 时会自动创建
		apiInfo.Domain = downloadDomain
		apiInfo.ToFile = filepath.Join(downloadCfg.DestDir, apiInfo.Key)
		apiInfo.StatusDBPath = dbPath
		apiInfo.Referer = downloadCfg.Referer
		apiInfo.FileEncoding = downloadCfg.FileEncoding
		apiInfo.Bucket = downloadCfg.Bucket
		if !downloadCfg.CheckHash {
			apiInfo.FileHash = ""
		}

		file, err := downloadFile(DownloadInfo{
			ApiInfo:  *apiInfo,
			IsPublic: downloadCfg.Public,
		})
		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(*download.ApiInfo)
		res := result.(download.ApiResult)
		export.Success().ExportF("download success, [%s:%s] => %s", apiInfo.Bucket, apiInfo.Key, res.FileAbsPath)
	}).OnWorkError(func(work work.Work, err error) {
		apiInfo := work.(*download.ApiInfo)
		export.Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* key fileSize fileHash and fileModifyTime */
			apiInfo.Key, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileHash, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()
}
