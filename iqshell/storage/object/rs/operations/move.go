package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type MoveInfo rs.MoveApiInfo

func Move(info MoveInfo) {
	result, err := rs.BatchOne(rs.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Move error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Move error:%v", result.Error)
		return
	}
}

type BatchMoveInfo struct {
	BatchInfo    BatchInfo
	SourceBucket string
	DestBucket   string
}

func BatchMove(info BatchMoveInfo) {
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

	rs.BatchWithHandler(&batchMoveHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchMoveHandler struct {
	scanner      *batchScanner
	info         *BatchMoveInfo
	resultExport *BatchResultExport
}

var _ rs.BatchHandler = (*batchMoveHandler)(nil)

func (b batchMoveHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchMoveHandler) ReadOperation() (rs.BatchOperation, bool) {
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
			info = &rs.MoveApiInfo{
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

func (b batchMoveHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.MoveApiInfo)
	if !ok {
		return
	}

	info := MoveInfo(apiInfo)
	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail().ExportF("%s\t%s\t%d\t%s\n", info.SourceKey, info.DestKey, result.Code, result.Error)
		log.ErrorF("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			result.Code, result.Error)
	} else {
		b.resultExport.Success().ExportF("%s\t%s\n", info.SourceKey, info.DestKey)
		log.ErrorF("Move '%s:%s' => '%s:%s' success\n",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey)
	}
}

func (b batchMoveHandler) HandlerError(err error) {
	log.ErrorF("batch move error:%v:", err)
}
