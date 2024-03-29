package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
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
	if err != nil || result == nil {
		data.SetCmdStatusError()
		log.ErrorF("Change mimetype Failed, [%s:%s] => '%s', Error:%v", info.Bucket, info.Key, info.Mime, err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Change mimetype Success, [%s:%s] => '%s'", info.Bucket, info.Key, info.Mime)
	} else {
		data.SetCmdStatusError()
		log.ErrorF("Change mimetype Failed, [%s:%s] => '%s', Code:%d, Error:%v",
			info.Bucket, info.Key, info.Mime, result.Code, result.Error)
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
			return &object.ChangeMimeApiInfo{}
		}).
		SetFileExport(exporter).
		ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
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
		}).
		OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
			apiInfo, ok := (operation).(*object.ChangeMimeApiInfo)
			if !ok {
				data.SetCmdStatusError()
				log.ErrorF("Change mimetype Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
				return
			}
			in := (*ChangeMimeInfo)(apiInfo)
			if result.IsSuccess() {
				log.InfoF("Change mimetype Success, [%s:%s] => '%s'", in.Bucket, in.Key, in.Mime)
			} else {
				data.SetCmdStatusError()
				log.ErrorF("Change mimetype Failed, [%s:%s] => '%s', Code: %d, Error: %s",
					in.Bucket, in.Key, in.Mime, result.Code, result.Error)
			}
		}).
		OnError(func(err *data.CodeError) {
			data.SetCmdStatusError()
			log.ErrorF("Batch change mimetype error:%v:", err)
		}).Start()
}
