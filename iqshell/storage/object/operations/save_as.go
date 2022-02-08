package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"os"
)

type SaveAsInfo object.SaveAsApiInfo

func SaveAs(info SaveAsInfo) {
	url, err := object.SaveAs(object.SaveAsApiInfo(info))
	if err != nil {
		log.ErrorF("save as error: %v", err)
		os.Exit(data.StatusError)
	} else {
		log.Alert(url)
	}
}
