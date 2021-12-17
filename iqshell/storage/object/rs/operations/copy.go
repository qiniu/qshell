package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type CopyInfo rs.CopyApiInfo

func Copy(info CopyInfo) {
	result, err := rs.BatchOne(rs.CopyApiInfo(info))
	if err != nil {
		log.ErrorF("Copy error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Copy error:%v", result.Error)
		return
	}
}

type BatchCopyInfo struct {
	BatchInfo    BatchInfo
	SourceBucket string
	DestBucket   string
}

func BatchCopy(info BatchCopyInfo) {
	if !prepareToBatch(info.BatchInfo) {
		return
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:     info.BatchInfo.FailExportFilePath,
		OverrideExportFilePath: info.BatchInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	scanner, err := newBatchScanner(info.BatchInfo)
	if err != nil {
		log.ErrorF("get scanner error:%v", err)
		return
	}

	rs.BatchWithHandler(&batchCopyHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchCopyHandler struct {
	scanner      *batchScanner
	info         *BatchCopyInfo
	resultExport *export.FileExporter
}

var _ rs.BatchHandler = (*batchCopyHandler)(nil)

func (b batchCopyHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchCopyHandler) ReadOperation() (rs.BatchOperation, bool) {
	var info rs.BatchOperation = nil

	line, success := b.scanner.scanLine()
	if !success {
		return nil, true
	}

	items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
	if len(items) > 0 {
		srcKey, destKey := items[0], items[0]
		if len(items) > 1 {
			destKey = items[1]
		}
		if srcKey != "" && destKey != "" {
			info = rs.CopyApiInfo{
				SourceBucket: b.info.SourceBucket,
				SourceKey:    srcKey,
				DestBucket:   b.info.DestBucket,
				DestKey:      destKey,
				Force:        b.info.BatchInfo.Force,
			}
		}
	}

	return info, false
}

func (b batchCopyHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.CopyApiInfo)
	if !ok {
		return
	}

	info := CopyInfo(apiInfo)
	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail().ExportF("%s\t%s\t%d\t%s\n", info.SourceKey, info.DestKey, result.Code, result.Error)
		log.ErrorF("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			result.Code, result.Error)
	} else {
		b.resultExport.Success().ExportF("%s\t%s\n", info.SourceKey, info.DestKey)
		log.ErrorF("Copy '%s:%s' => '%s:%s' success\n",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey)
	}
}

func (b batchCopyHandler) HandlerError(err error) {
	log.ErrorF("batch copy error:%v:", err)
}
