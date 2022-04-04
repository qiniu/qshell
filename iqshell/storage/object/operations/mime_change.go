package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type ChangeMimeInfo object.ChangeMimeApiInfo

func (info *ChangeMimeInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.Mime) == 0 {
		return alert.CannotEmptyError("MimeType", "")
	}
	return nil
}

func ChangeMime(cfg *iqshell.Config, info ChangeMimeInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.ChangeMimeType((*object.ChangeMimeApiInfo)(&info))
	if err != nil {
		log.ErrorF("Change mimetype Failed, [%s:%s] => '%d', Error:%v", info.Bucket, info.Key, info.Mime, err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change mimetype Failed, [%s:%s] => '%s', Code:%d, Error:%v", info.Bucket, info.Key, info.Mime, result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Change mimetype Success, [%s:%s] => '%s'", info.Bucket, info.Key, info.Mime)
	}
}

type BatchChangeMimeInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchChangeMimeInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchChangeMime(cfg *iqshell.Config, info BatchChangeMimeInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}
	batch.NewHandler(info.BatchInfo).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		if len(items) > 1 {
			key, mime := items[0], items[1]
			if key != "" && mime != "" {
				return &object.ChangeMimeApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Mime:   mime,
				}, nil
			}
			return nil, alert.Error("key or mime invalid", "")
		}
		return nil, alert.Error("need more than one param", "")
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.ChangeMimeApiInfo)
		if !ok {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Change mimetype Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
			return
		}
		in := (*ChangeMimeInfo)(apiInfo)
		if result.Code != 200 || result.Error != "" {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Change mimetype Failed, [%s:%s] => '%s', Code: %d, Error: %s",
				in.Bucket, in.Key, in.Mime, result.Code, result.Error)
		} else {
			exporter.Success().Export(operationInfo)
			log.InfoF("Change mimetype Success, [%s:%s] => '%s'", in.Bucket, in.Key, in.Mime)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch change mimetype error:%v:", err)
	}).Start()
}
