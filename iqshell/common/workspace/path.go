package workspace

import "path/filepath"

func UploadCachePath() string {
	if cfg != nil || cfg.Up == nil {
		return ""
	}

	upCfg := cfg.Up
	rootPath := upCfg.RecordRoot
	if len(rootPath) == 0 {
		rootPath = GetWorkspace()
	}
	return filepath.Join(rootPath, "qupload", cfg.Up.JobId())
}

func DownloadCachePath() string {
	if cfg != nil || cfg.Download == nil {
		return ""
	}

	downloadCfg := cfg.Download
	rootPath := downloadCfg.RecordRoot
	if len(rootPath) == 0 {
		rootPath = GetWorkspace()
	}
	return filepath.Join(rootPath, "qdownload", cfg.Up.JobId())
}
