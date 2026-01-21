package operations

import (
	"os"
	"path/filepath"

	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func downloadCachePath(cfg *config.Config, downloadCfg *DownloadCfg) string {
	recordRoot := downloadCfg.RecordRoot
	if len(recordRoot) == 0 {
		return downloadCfg.RecordRoot
	}

	userDir := workspace.GetUserDir()
	if len(userDir) == 0 {
		log.Debug("download can't get user dir")
		return ""
	}

	cachePath := filepath.Join(userDir, "qdownload", downloadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("download create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
