package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path/filepath"
)

func uploadCachePath(cfg *config.Config, uploadCfg *UploadConfig) string {
	recordRoot := uploadCfg.RecordRoot
	if len(recordRoot) == 0 {
		recordRoot = workspace.GetUserPath()
	}

	if len(recordRoot) == 0 {
		log.Debug("upload can't get record root")
		return ""
	}

	cachePath := filepath.Join(recordRoot, "qupload", uploadCfg.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("upload create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
