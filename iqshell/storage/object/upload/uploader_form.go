package upload

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
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

func (f *formUploader) upload(info ApiInfo) (ret Result, err error) {
	file, err := os.Open(info.FilePath)
	if err != nil {
		err =  errors.New("form upload: open file error:" + err.Error())
		return
	}

	fileStatus, err := file.Stat()
	if err != nil {
		err =  errors.New("form upload: ger file status error:" + err.Error())
		return
	}

	up := storage.NewFormUploader(f.cfg)
	err = up.Put(workspace.GetContext(), &ret, info.TokenProvider(), info.SaveKey, file, fileStatus.Size(), f.ext)
	if err != nil {
		err =  errors.New("form upload: upload error:" + err.Error())
	}

	return
}

