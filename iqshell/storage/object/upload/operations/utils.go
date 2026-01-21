package operations

import (
	"os"
	"path/filepath"

	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func uploadCachePath(cfg *config.Config, uploadCfg *UploadConfig) string {
	recordRoot := uploadCfg.RecordRoot
	if len(recordRoot) == 0 {
		return uploadCfg.RecordRoot
	}

	userDir := workspace.GetUserDir()
	if len(userDir) == 0 {
		log.Debug("upload can't get user dir")
		return ""
	}

	cachePath := filepath.Join(userDir, "qupload", uploadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("upload create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
