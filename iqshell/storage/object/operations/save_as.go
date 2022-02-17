package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"os"
)

type SaveAsInfo object.SaveAsApiInfo

func (info *SaveAsInfo) Check() error {
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

func SaveAs(info SaveAsInfo) {
	url, err := object.SaveAs(object.SaveAsApiInfo(info))
	if err != nil {
		log.ErrorF("save as error: %v", err)
		os.Exit(data.StatusError)
	} else {
		log.Alert(url)
	}
}
