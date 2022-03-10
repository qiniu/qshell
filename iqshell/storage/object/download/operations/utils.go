package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"path/filepath"
)

func downloadCachePath(cfg *config.Config, downloadCfg *DownloadCfg) string {
	recordRoot := downloadCfg.RecordRoot
	if len(recordRoot) == 0 {
		if cfg == nil {
			return ""
		}
		recordRoot = cfg.RecordRoot.Value()
	}

	cachePath := filepath.Join(recordRoot, "qdownload", downloadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("download create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
