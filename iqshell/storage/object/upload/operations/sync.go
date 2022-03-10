package operations

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type SyncInfo UploadInfo

func (info *SyncInfo) Check() error {
	if len(info.FilePath) == 0 {
		return alert.CannotEmptyError("SrcResUrl", "")
	}
	if !utils.IsNetworkSource(info.FilePath) {
		return alert.Error("sync only for network source", "")
	}
	if len(info.ToBucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func SyncFile(cfg *iqshell.Config, info SyncInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	info.Progress = progress.NewPrintProgress(" 进度")
	ret, err := uploadFile((*UploadInfo)(&info))
	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			log.ErrorF("Sync file error %d: %s, Reqid: %s", v.Code, v.Err, v.Reqid)
		}
	} else {
		log.Alert("")
		log.Alert("-------------- File Info --------------")
		log.AlertF("%10s%s", "Key: ", ret.Key)
		log.AlertF("%10s%s", "Hash: ", ret.Hash)
		log.AlertF("%10s%d%s", "Fsize: ", ret.FSize, "("+utils.FormatFileSize(ret.FSize)+")")
		log.AlertF("%10s%s", "MimeType: ", ret.MimeType)
	}
}
