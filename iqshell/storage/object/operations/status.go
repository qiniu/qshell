package operations

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type StatusInfo object.StatusApiInfo

func (info *StatusInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func Status(cfg *iqshell.Config, info StatusInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Status(object.StatusApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Status Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if result.IsSuccess() {
		log.Alert(getResultInfo(info.Bucket, info.Key, result))
	}
}

type BatchStatusInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchStatusInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchStatus(cfg *iqshell.Config, info BatchStatusInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", cfg.CmdCfg.CmdId, info.Bucket, info.BatchInfo.InputFile))
		return filepath.Join(cmdPath, jobId)
	}
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
		return
	}

	batch.NewHandler(info.BatchInfo).
		EmptyOperation(func() flow.Work {
			return &object.StatusApiInfo{}
		}).
		SetFileExport(exporter).
		ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
			key := items[0]
			if key != "" {
				return &object.StatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
				}, nil
			}
			return nil, alert.Error("key invalid", "")
		}).
		OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
			apiInfo, ok := (operation).(*object.StatusApiInfo)
			if !ok {
				data.SetCmdStatusError()
				log.ErrorF("Status Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
				return
			}
			in := (*StatusInfo)(apiInfo)
			if result.IsSuccess() {
				infoString := fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d",
					in.Key, result.FSize, result.Hash, result.MimeType, result.PutTime, result.Type)
				log.Alert(infoString)
				exporter.Result().Export(infoString)
			} else {
				data.SetCmdStatusError()
				log.ErrorF("Status Failed, [%s:%s], Code: %d, Error: %s", in.Bucket, in.Key, result.Code, result.Error)
			}
		}).
		OnError(func(err *data.CodeError) {
			data.SetCmdStatusError()
			log.ErrorF("Batch Status error:%v:", err)
		}).Start()
}

func getResultInfo(bucket, key string, status object.StatusResult) string {
	statInfo := ""
	fieldAdder := func(name string, value interface{}, desc string) {
		if len(desc) == 0 {
			statInfo += fmt.Sprintf("%-25s%v\r\n", name+":", value)
		} else {
			statInfo += fmt.Sprintf("%-25s%v -> %s\r\n", name+":", value, desc)
		}
	}

	fieldAdderWithValueDescs := func(name string, value interface{}, descs map[interface{}]string, noneDesc string) {
		desc := descs[value]
		if len(descs[value]) > 0 {
			fieldAdder(name, value, desc)
		} else {
			fieldAdder(name, noneDesc, "")
		}
	}

	fieldAdder("Bucket", bucket, "")
	fieldAdder("Key", key, "")
	fieldAdder("Etag", status.Hash, "")
	fieldAdder("MD5", status.MD5, "")
	fieldAdder("Fsize", status.FSize, utils.FormatFileSize(status.FSize))
	fieldAdder("PutTime", status.PutTime, time.Unix(0, status.PutTime*100).String())
	fieldAdder("MimeType", status.MimeType, "")

	fieldAdderWithValueDescs("Status", status.Status,
		map[interface{}]string{1: "禁用"},
		"未禁用")

	fieldAdderWithValueDescs("RestoreStatus", status.Status,
		map[interface{}]string{1: "解冻中", 2: "解冻完成"},
		"无解冻操作")

	lifecycleFieldAdder := func(name string, date int64) {
		if date > 0 {
			t := time.Unix(date, 0)
			fieldAdder(name, date, t.String())
		} else {
			fieldAdder(name, "未设置", "")
		}
	}
	lifecycleFieldAdder("Expiration", status.Expiration)
	lifecycleFieldAdder("TransitionToIA", status.TransitionToIA)
	lifecycleFieldAdder("TransitionToArchive", status.TransitionToARCHIVE)
	lifecycleFieldAdder("TransitionToDeepArchive", status.TransitionToDeepArchive)

	fieldAdder("FileType", status.Type, getFileTypeDescription(status.Type))

	return statInfo
}

var objectTypes = []string{"标准存储", "低频存储", "归档存储", "深度归档存储"}

func getFileTypeDescription(fileTypes int) string {
	typeString := "未知类型"
	if fileTypes >= 0 && fileTypes < len(objectTypes) {
		typeString = objectTypes[fileTypes]
	}
	return typeString
}
