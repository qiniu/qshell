package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type ChangeMimeInfo rs.ChangeMimeApiInfo

func ChangeMime(info ChangeMimeInfo) {
	result, err := rs.BatchOne(rs.ChangeMimeApiInfo(info))
	if err != nil {
		log.ErrorF("Change Mime error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Mime:%v", result.Error)
		return
	}
}

type BatchChangeMimeInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

func BatchChangeMime(info BatchChangeMimeInfo) {
	if !prepareToBatch(info.BatchInfo) {
		return
	}

	resultExport, err := NewBatchResultExport(info.BatchInfo)
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	scanner, err := newBatchScanner(info.BatchInfo)
	if err != nil {
		log.ErrorF("get scanner error:%v", err)
		return
	}

	rs.BatchWithHandler(&batchChangeMimeHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchChangeMimeHandler struct {
	scanner      *batchScanner
	info         *BatchChangeMimeInfo
	resultExport *BatchResultExport
}

var _ rs.BatchHandler = (*batchChangeMimeHandler)(nil)

func (b batchChangeMimeHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchChangeMimeHandler) ReadOperation() (rs.BatchOperation, bool) {
	var info rs.BatchOperation = nil

	line, success := b.scanner.scanLine()
	if !success {
		return nil, true
	}

	items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
	if len(items) > 1 {
		key, mime := items[0], items[1]
		if key != "" && mime != "" {
			info = rs.ChangeMimeApiInfo{
				Bucket: b.info.Bucket,
				Key:    key,
				Mime:   mime,
			}
		}
	}

	return info, false
}

func (b batchChangeMimeHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.ChangeMimeApiInfo)
	if !ok {
		return
	}

	info := ChangeMimeInfo(apiInfo)
	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%s\t%d\t%s\n", info.Key, info.Mime, result.Code, result.Error)
		log.ErrorF("Chgm '%s' => '%s' Failed, Code: %d, Error: %s\n",
			info.Key, info.Mime, result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%s\n", info.Key, info.Mime)
		log.ErrorF("Chgm '%s' => '%s' success\n", info.Key, info.Mime)
	}
}

func (b batchChangeMimeHandler) HandlerError(err error) {
	log.ErrorF("batch chgm error:%v:", err)
}
