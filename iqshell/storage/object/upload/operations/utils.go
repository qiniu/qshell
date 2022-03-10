package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"path/filepath"
)

func uploadCachePath(cfg *config.Config, uploadCfg *UploadConfig) string {
	recordRoot := uploadCfg.RecordRoot
	if len(recordRoot) == 0 {
		if cfg == nil {
			return ""
		}
		recordRoot = cfg.RecordRoot.Value()
	}

	cachePath := filepath.Join(recordRoot, "qupload", uploadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("upload create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
