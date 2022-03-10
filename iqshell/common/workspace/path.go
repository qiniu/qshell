package workspace

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"path/filepath"
)

func UploadCachePath() string {
	if cfg == nil || cfg.Up == nil {
		return ""
	}

	upCfg := cfg.Up
	rootPath := GetWorkspace()
	if data.NotEmpty(upCfg.RecordRoot) {
		rootPath = upCfg.RecordRoot.Value()
	}

	cachePath := filepath.Join(rootPath, "qupload", cfg.Up.JobId())
	if cErr := os.MkdirAll(cachePath, os.ModePerm); cErr != nil {
		log.WarningF("upload create cache dir error:%v", cErr)
		return ""
	}
	return cachePath
}
