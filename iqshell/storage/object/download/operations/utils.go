package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path/filepath"
)

func downloadCachePath(cfg *config.Config, downloadCfg *DownloadCfg) string {
	if len(downloadCfg.RecordRoot) > 0 {
		return downloadCfg.RecordRoot
	}

	userPath := workspace.GetUserPath()
	if len(userPath) == 0 {
		log.Debug("download can't get user dir")
		return ""
	}

	cachePath := filepath.Join(userPath, "qdownload", downloadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("download create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
