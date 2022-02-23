package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

type DirCacheInfo struct {
	Dir        string
	SaveToFile string
}

func (info *DirCacheInfo) Check() error {
	if len(info.Dir) == 0 {
		return alert.CannotEmptyError("directory path", "")
	}
	return nil
}

func DirCache(cfg *iqshell.Config, info DirCacheInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
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
