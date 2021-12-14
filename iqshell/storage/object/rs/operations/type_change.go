package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
)

type ChangeTypeInfo struct {
	Bucket string
	Key    string
	Type   string
}

func (c ChangeTypeInfo) getTypeOfInt() (int, error) {
	if len(c.Type) == 0 {
		return -1, errors.New(alert.CannotEmpty("type", ""))
	}

	ret, err := strconv.Atoi(c.Type)
	if err != nil {
		return -1, errors.New("parse type error:" + err.Error())
	}

	if ret < 0 || ret > 1 {
		return -1, errors.New("type must be 0 or 1")
	}
	return ret, nil
}

func ChangeType(info ChangeTypeInfo) {
	t, err := info.getTypeOfInt()
	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	result, err := rs.ChangeType(rs.ChangeTypeApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Type:   t,
	})

	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Type:%v", result.Error)
		return
	}
}

type BatchChangeTypeInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

func BatchChangeType(info BatchChangeTypeInfo) {
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

	rs.BatchWithHandler(&batchChangeTypeHandler{
		scanner:      scanner,
		info:         &info,
		resultExport: resultExport,
	})
}

type batchChangeTypeHandler struct {
	scanner      *batchScanner
	info         *BatchChangeTypeInfo
	resultExport *BatchResultExport
}

var _ rs.BatchHandler = (*batchChangeTypeHandler)(nil)

func (b batchChangeTypeHandler) WorkCount() int {
	return b.info.BatchInfo.Worker
}

func (b batchChangeTypeHandler) ReadOperation() rs.BatchOperation {
	var info *rs.ChangeTypeApiInfo

	for {
		line, complete := b.scanner.scanLine()
		if complete {
			break
		}

		items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
		if len(items) > 1 {
			key, t := items[0], items[1]
			tInt, err := strconv.Atoi(t)
			if err != nil {
				continue
			}

			if key != "" && t != "" {
				info = &rs.ChangeTypeApiInfo{
					Bucket: b.info.Bucket,
					Key:    key,
					Type:   tInt,
				}
				break
			}
		}
	}

	return info
}

func (b batchChangeTypeHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
	info, ok := (operation).(rs.ChangeTypeApiInfo)
	if !ok {
		return
	}

	if result.Code != 200 || result.Error != "" {
		b.resultExport.Fail.ExportF("%s\t%d\t%d\t%s\n", info.Key, info.Type, result.Code, result.Error)
		log.ErrorF("Change Type '%s' => '%s' Failed, Code: %d, Error: %s\n",
			info.Key, info.Type, result.Code, result.Error)
	} else {
		b.resultExport.Success.ExportF("%s\t%d\n", info.Key, info.Type)
		log.ErrorF("Change Type '%s' => '%d' success\n", info.Key, info.Type)
	}
}

func (b batchChangeTypeHandler) HandlerError(err error) {
	log.ErrorF("batch change Type error:%v:", err)
}
