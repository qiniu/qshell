package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type CopyInfo object.CopyApiInfo

func (info *CopyInfo) Check() error {
	if len(info.SourceBucket) == 0 {
		return alert.CannotEmptyError("SourceBucket", "")
	}
	if len(info.SourceKey) == 0 {
		return alert.CannotEmptyError("SourceKey", "")
	}
	if len(info.DestBucket) == 0 {
		return alert.CannotEmptyError("DestBucket", "")
	}
	if len(info.DestKey) == 0 {
		return alert.CannotEmptyError("DestKey", "")
	}
	return nil
}

func Copy(info CopyInfo) {
	result, err := object.Copy(object.CopyApiInfo(info))
	if err != nil {
		log.ErrorF("Copy error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Copy error:%v", result.Error)
		return
	}

	log.InfoF("copy success [%s:%s] => [%s:%s]",
		info.SourceBucket, info.SourceKey,
		info.DestBucket, info.DestKey)
}

type BatchCopyInfo struct {
	BatchInfo    batch.Info
	SourceBucket string
	DestBucket   string
}

func (info *BatchCopyInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.SourceBucket) == 0 {
		return alert.CannotEmptyError("SrcBucket", "")
	}

	if len(info.DestBucket) == 0 {
		return alert.CannotEmptyError("DestBucket", "")
	}

	return nil
}

func BatchCopy(info BatchCopyInfo) {
	handler, err := group.NewHandler(info.BatchInfo.Info)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewFlow(info.BatchInfo).ReadOperation(func() (operation batch.Operation, hasMore bool) {
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, false
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) > 0 {
			// 如果只有一个参数，源 key 即为目标 key
			srcKey, destKey := items[0], items[0]
			if len(items) > 1 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				return object.CopyApiInfo{
					SourceBucket: info.SourceBucket,
					SourceKey:    srcKey,
					DestBucket:   info.DestBucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Overwrite,
				}, true
			} else {
				return nil, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.CopyApiInfo)
		if !ok {
			return
		}
		in := CopyInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s", in.SourceKey, in.DestKey, result.Code, result.Error)
			log.ErrorF("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s", in.SourceKey, in.DestKey)
			log.InfoF("Copy '%s:%s' => '%s:%s' success",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch copy error:%v:", err)
	}).Start()
}
