package upload

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
)

type formUploader struct {
	cfg *storage.Config
	ext *storage.PutExtra
}

func newFromUploader(cfg *storage.Config, ext *storage.PutExtra) Uploader {
	return &formUploader{
		cfg: cfg,
		ext: ext,
	}
}

func (f *formUploader) upload(info ApiInfo) (ret ApiResult, err error) {
	log.DebugF("form upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)

	file, err := os.Open(info.FilePath)
	if err != nil {
		err = errors.New("form upload: open file error:" + err.Error())
		return
	}

	fileStatus, err := file.Stat()
	if err != nil {
		err = errors.New("form upload: ger file status error:" + err.Error())
		return
	}

	token := info.TokenProvider()
	log.DebugF("upload token:%s", token)

	if info.Progress != nil {
		info.Progress.SetFileSize(info.FileSize)
		info.Progress.Start()
		f.ext.OnProgress = func(fsize, uploaded int64) {
			info.Progress.Progress(uploaded)
		}
	}

	up := storage.NewFormUploader(f.cfg)
	err = up.Put(workspace.GetContext(), &ret, token, info.SaveKey, file, fileStatus.Size(), f.ext)
	if err != nil {
		err = errors.New("form upload: upload error:" + err.Error())
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
	}

	return
}
