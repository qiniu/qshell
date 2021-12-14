package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
)

type ForbiddenInfo struct {
	Bucket      string
	Key         string
	UnForbidden bool
}

func (c ForbiddenInfo) getStatus() int {
	// 0:启用  1:禁用
	if c.UnForbidden {
		return 0
	} else {
		return 1
	}
}

func ForbiddenObject(info ForbiddenInfo) {
	result, err := rs.ChangeStatus(rs.ChangeStatusApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Status: info.getStatus(),
	})

	if err != nil {
		log.ErrorF("change stat error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("change stat error:%s", result.Error)
		return
	}
}

type BatchChangeStatusInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

func BatchChangeStatus(info BatchChangeStatusInfo) {
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

	rs.BatchWithHandler(&batchChangeStatusHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchChangeStatusHandler struct {
	scanner      *batchScanner
	info         *BatchChangeStatusInfo
	resultExport *BatchResultExport
}

var _ rs.BatchHandler = (*batchChangeStatusHandler)(nil)

func (b batchChangeStatusHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchChangeStatusHandler) ReadOperation() rs.BatchOperation {
	var info *rs.ChangeStatusApiInfo

	for {
		line, complete := b.scanner.scanLine()
		if complete {
			break
		}

		items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
		if len(items) > 1 {
			key, status := items[0], items[1]
			statusInt, err := strconv.Atoi(status)
			if err != nil {
				continue
			}

			if key != "" && status != "" {
				info = &rs.ChangeStatusApiInfo{
					Bucket: b.info.Bucket,
					Key:    key,
					Status: statusInt,
				}
				break
			}
		}
	}

	return info
}

func (b batchChangeStatusHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	info, ok := (operation).(rs.ChangeStatusApiInfo)
	if !ok {
		return
	}

	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%d\t%d\t%s\n", info.Key, info.Status, result.Code, result.Error)
		log.ErrorF("Change status '%s' => '%s' Failed, Code: %d, Error: %s\n",
			info.Key, info.Status, result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%d\n", info.Key, info.Status)
		log.ErrorF("Change status '%s' => '%d' success\n", info.Key, info.Status)
	}
}

func (b batchChangeStatusHandler) HandlerError(err error) {
	log.ErrorF("batch change status error:%v:", err)
}
