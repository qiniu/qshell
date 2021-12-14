package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

// rename 实际用的还是 move

type RenameInfo rs.MoveApiInfo

func Rename(info RenameInfo) {
	result, err := rs.Move(rs.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Rename error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Rename error:%v", result.Error)
		return
	}
}

type BatchRenameInfo struct {
	BatchInfo    BatchInfo
	SourceBucket string
	DestBucket   string
}

func BatchRename(info BatchRenameInfo) {
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

	rs.BatchWithHandler(&batchRenameHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchRenameHandler struct {
	scanner      *batchScanner
	info         *BatchRenameInfo
	resultExport *BatchResultExport
}

var _ rs.BatchHandler = (*batchRenameHandler)(nil)

func (b batchRenameHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchRenameHandler) ReadOperation() rs.BatchOperation {
	var info *rs.MoveApiInfo

	for {
		line, complete := b.scanner.scanLine()
		if complete {
			break
		}

		items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
		if len(items) > 1 {
			sourceKey, destKey := items[0], items[1]
			if sourceKey != "" && destKey != "" {
				info = &rs.MoveApiInfo{
					SourceBucket: b.info.SourceBucket,
					SourceKey:    sourceKey,
					DestBucket:   b.info.DestBucket,
					DestKey:      destKey,
					Force:        b.info.BatchInfo.Force,
				}
				break
			}
		}
	}

	return info
}

func (b batchRenameHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.MoveApiInfo)
	if !ok {
		return
	}

	info := RenameInfo(apiInfo)
	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%s\t%d\t%s\n", info.SourceKey, info.DestKey, result.Code, result.Error)
		log.ErrorF("Rename '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%s\n", info.SourceKey, info.DestKey)
		log.ErrorF("Rename '%s:%s' => '%s:%s' success\n",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey)
	}
}

func (b batchRenameHandler) HandlerError(err error) {
	log.ErrorF("batch rename error:%v:", err)
}
