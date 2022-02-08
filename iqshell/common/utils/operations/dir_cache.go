package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

type DirCacheInfo struct {
	Dir        string
	SaveToFile string
}

func DirCache(info DirCacheInfo) {
	if len(info.Dir) == 0 {
		log.Error(alert.CannotEmpty("directory path", ""))
		return
	}

	if info.SaveToFile == "" {
		info.SaveToFile = "stdout"
	}

	_, retErr := utils.DirCache(info.Dir, info.SaveToFile)
	if retErr != nil {
		os.Exit(data.StatusError)
	}
}
