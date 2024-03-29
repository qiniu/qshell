package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
)

type SaveAsInfo object.SaveAsApiInfo

func (info *SaveAsInfo) Check() *data.CodeError {
	if len(info.PublicUrl) == 0 {
		return alert.CannotEmptyError("PublicUrlWithFop", "")
	}
	if len(info.SaveBucket) == 0 {
		return alert.CannotEmptyError("SaveBucket", "")
	}
	if len(info.SaveKey) == 0 {
		return alert.CannotEmptyError("SaveKey", "")
	}
	return nil
}

func SaveAs(cfg *iqshell.Config, info SaveAsInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	url, err := object.SaveAs(object.SaveAsApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Save as Failed, Error: %v", err)
	} else {
		log.Alert(url)
	}
}
