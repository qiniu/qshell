package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
)

type DeleteInfo struct {
	Bucket    string
	Key       string
	AfterDays string
}

func getAfterDaysOfInt(after string) (int, error) {
	if len(after) == 0 {
		return -1, nil
	}
	return strconv.Atoi(after)
}

func Delete(info DeleteInfo) {
	afterDays, err := getAfterDaysOfInt(info.AfterDays)
	if err != nil {
		log.ErrorF("delete after days invalid:%v", err)
		return
	}

	result, err := rs.BatchOne(rs.DeleteApiInfo{
		Bucket:    info.Bucket,
		Key:       info.Key,
		AfterDays: afterDays,
	})

	if err != nil {
		log.ErrorF("delete error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("delete error:%s", result.Error)
		return
	}
}

type BatchDeleteInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

// BatchDelete 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDelete(info BatchDeleteInfo) {
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

	rs.BatchWithHandler(&batchDeleteHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})

}

type batchDeleteHandler struct {
	scanner      *batchScanner
	info         *BatchDeleteInfo
	resultExport *BatchResultExport
}

func (b batchDeleteHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchDeleteHandler) ReadOperation() (rs.BatchOperation, bool) {
	var info rs.BatchOperation = nil

	line, success := b.scanner.scanLine()
	if !success {
		return nil, true
	}

	items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
	if len(items) > 0 {
		key := items[0]
		putTime := ""
		if len(items) > 1 {
			putTime = items[1]
		}

		if key != "" {
			info = rs.DeleteApiInfo{
				Bucket: b.info.Bucket,
				Key:    key,
				Condition: rs.OperationCondition{
					PutTime: putTime,
				},
			}
		}
	}

	return info, false
}

func (b batchDeleteHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.DeleteApiInfo)
	if !ok {
		return
	}

	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%s\t%d\t%s\n", apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
		log.ErrorF("Delete '%s' when put time:'%d' Failed, Code: %d, Error: %s\n", apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%s\n", apiInfo.Key, apiInfo.Condition.PutTime)
		log.AlertF("Delete '%s' when put time:'%d' success\n", apiInfo.Key, apiInfo.Condition.PutTime)
	}
}

func (b batchDeleteHandler) HandlerError(err error) {
	log.ErrorF("batch delete error:%v:", err)
}

// BatchDeleteAfter 延迟批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDeleteAfter(info BatchDeleteInfo) {
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

	rs.BatchWithHandler(&batchDeleteAfterHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})

}

type batchDeleteAfterHandler struct {
	scanner      *batchScanner
	info         *BatchDeleteInfo
	resultExport *BatchResultExport
}

func (b batchDeleteAfterHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchDeleteAfterHandler) ReadOperation() (rs.BatchOperation, bool) {
	var info rs.BatchOperation = nil

	line, success := b.scanner.scanLine()
	if !success {
		return nil, true
	}

	items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
	if len(items) > 0 {
		key := items[0]
		after := ""
		if len(items) > 1 {
			after = items[1]
		}

		afterDays, err := getAfterDaysOfInt(after)
		if err != nil {
			log.ErrorF("parse after days error:%v from:%s", err, after)
		} else if key != "" {
			info = rs.DeleteApiInfo{
				Bucket:    b.info.Bucket,
				Key:       key,
				AfterDays: afterDays,
			}
		}
	}

	return info, false
}

func (b batchDeleteAfterHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	apiInfo, ok := (operation).(rs.DeleteApiInfo)
	if !ok {
		return
	}

	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%s\t%d\t%s\n", apiInfo.Key, apiInfo.AfterDays, result.Code, result.Error)
		log.ErrorF("Expire '%s' => '%d' Failed, Code: %d, Error: %s\n", apiInfo.Key, apiInfo.AfterDays, result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%s\n", apiInfo.Key, apiInfo.AfterDays)
		log.AlertF("Expire '%s' => '%d' success\n", apiInfo.Key, apiInfo.AfterDays)
	}
}

func (b batchDeleteAfterHandler) HandlerError(err error) {
	log.ErrorF("batch delete after error:%v:", err)
}
