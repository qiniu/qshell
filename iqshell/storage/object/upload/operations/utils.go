package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path/filepath"
)

func uploadCachePath(cfg *config.Config, uploadCfg *UploadConfig) string {
	if len(uploadCfg.RecordRoot) > 0 {
		return uploadCfg.RecordRoot
	}

	userPath := workspace.GetUserPath()
	if len(userPath) == 0 {
		log.Debug("upload can't get user dir")
		return ""
	}

	cachePath := filepath.Join(userPath, "qupload", uploadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("upload create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
