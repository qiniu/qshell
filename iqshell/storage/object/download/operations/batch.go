package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/scanner"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
	"strconv"
)

type BatchDownloadInfo struct {
	GroupInfo group.Info
}

func BatchDownload(info BatchDownloadInfo) {
	downloadCfg := workspace.GetConfig().Download

	export, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.GroupInfo.SuccessExportFilePath,
		FailExportFilePath:     info.GroupInfo.FailExportFilePath,
		OverrideExportFilePath: info.GroupInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.Error(err)
		return
	}

	reader, err := newDownloadWorkReader(info.GroupInfo.InputFile, info.GroupInfo.ItemSeparate, downloadCfg, export)
	if err != nil {
		log.Error(err)
		return
	}

	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		return reader.read()
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
		export.Success().ExportF("download success, [%s:%s] => %s", apiInfo.Bucket, apiInfo.Key, file)
	}).OnWorkError(func(work work.Work, err error) {
		apiInfo := work.(download.ApiInfo)
		export.Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* key fileSize fileHash and fileModifyTime */
			apiInfo.Key, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileHash, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()
}

type downloadWorkReader struct {
	exporter       *export.FileExporter
	lineScanner    scanner.Scanner
	itemSeparate   string
	inputFile      string
	infoChan       chan download.ApiInfo
	downloadCfg    *config.Download
	downloadDomain string
	dbDir          string
}

func newDownloadWorkReader(inputFile string, itemSeparate string, downloadCfg *config.Download, exporter *export.FileExporter) (r *downloadWorkReader, err error) {
	r = &downloadWorkReader{
		exporter:     exporter,
		itemSeparate: itemSeparate,
		downloadCfg:  downloadCfg,
		inputFile:    inputFile,
		infoChan:     make(chan download.ApiInfo, 100),
	}

	r.downloadDomain = downloadCfg.DownloadDomain()
	if len(r.downloadDomain) == 0 {
		r.downloadDomain, _ = bucket.DomainOfBucket(downloadCfg.Bucket)
	}
	if len(r.downloadDomain) == 0 {
		log.Error("download domain can't be empty, you can set cdn_domain or io_host")
		return
	}

	if len(downloadCfg.RecordRoot) == 0 {
		r.dbDir = filepath.Join(workspace.GetWorkspace(), "download")
	} else {
		r.dbDir = filepath.Join(downloadCfg.RecordRoot, "download")
	}
	log.InfoF("download db dir:%s", r.dbDir)

	r.lineScanner, err = scanner.NewScanner(scanner.Info{
		StdInEnable: true,
		InputFile:   inputFile,
	})

	r.createReadOperation()

	return
}

func (d *downloadWorkReader) createReadOperation() {
	go func() {
		var keys []string
		for {
			if len(keys) == 100 {
				d.statusAndAddToChan(keys)
				keys = nil
			}
			if keys == nil {
				keys = make([]string, 0, 100)
			}

			line, success := d.lineScanner.ScanLine()
			if !success {
				close(d.infoChan)
				break
			}

			items := utils.SplitString(line, d.itemSeparate)
			if len(items) < 4 {
				log.ErrorF("invalid line, line must contain key fileSize fileHash and fileModifyTime:%s", line)
				d.exporter.Fail().ExportF("%s: error:%s", line, "line must contain key fileSize fileHash and fileModifyTime")
				continue
			}

			fileKey := items[0]
			fileSize, err := strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				log.ErrorF("invalid line, get file size error:%s", line)
				d.exporter.Fail().ExportF("%s: get file size error:%s", line, err)
				continue
			}

			fileHash := items[2]
			fileModifyTime, err := strconv.ParseInt(items[3], 10, 64)
			if err != nil {
				log.ErrorF("invalid line, get file modify time error:%s", line)
				d.exporter.Fail().ExportF("%s: get file modify time error:%s", line, err)
				continue
			}

			d.infoChan <- download.ApiInfo{
				Url:            "", // downloadFile 时会自动创建
				Domain:         d.downloadDomain,
				ToFile:         filepath.Join(d.downloadCfg.DestDir, fileKey),
				StatusDBPath:   d.dbDir,
				Referer:        d.downloadCfg.Referer,
				FileEncoding:   d.downloadCfg.FileEncoding,
				Bucket:         d.downloadCfg.Bucket,
				Key:            fileKey,
				FileHash:       fileHash,
				FileSize:       fileSize,
				FileModifyTime: fileModifyTime,
			}
		}
	}()
}

func (d *downloadWorkReader) statusAndAddToChan(keys []string) {
	operations := make([]batch.Operation, 0, len(keys))
	for _, key := range keys {
		if len(key) > 0 {
			operations = append(operations, object.StatusApiInfo{
				Bucket: d.downloadCfg.Bucket,
				Key:    key,
			})
		}
	}
	results, err := batch.Some(operations)
	if err != nil {
		log.ErrorF("happen error:%v", err)
	}

	if len(results) == len(operations) {
		for i, result := range results {
			item := operations[i].(object.StatusApiInfo)
			if result.Code != 200 || result.Error != "" {
				d.exporter.Fail().ExportF("%s%s error:%v", item.Key, result.Error)
			} else {
				d.infoChan <- download.ApiInfo{
					Url:            "", // downloadFile 时会自动创建
					Domain:         d.downloadDomain,
					ToFile:         filepath.Join(d.downloadCfg.DestDir, item.Key),
					StatusDBPath:   d.dbDir,
					Referer:        d.downloadCfg.Referer,
					FileEncoding:   d.downloadCfg.FileEncoding,
					Bucket:         d.downloadCfg.Bucket,
					Key:            item.Key,
					FileHash:       result.Hash,
					FileSize:       result.FSize,
					FileModifyTime: result.PutTime,
				}
			}
		}
	}
}

func (d *downloadWorkReader) read() (info download.ApiInfo, hasMore bool) {
	for info = range d.infoChan {
		hasMore = true
	}
	return
}
