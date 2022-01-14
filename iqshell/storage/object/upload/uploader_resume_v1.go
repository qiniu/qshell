package upload

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
)

type resumeV1Uploader struct {
	cfg *storage.Config
	ext *storage.RputExtra
}

func newResumeV1Uploader(cfg *storage.Config, ext *storage.RputExtra) Uploader {
	return &resumeV1Uploader{
		cfg: cfg,
		ext: ext,
	}
}

func (r *resumeV1Uploader) upload(info ApiInfo) (ret UploadResult, err error) {
	file, err := os.Open(info.FilePath)
	if err != nil {
		err = errors.New("resume v1 upload: open file error:" + err.Error())
		return
	}

	fileStatus, err := file.Stat()
	if err != nil {
		err = errors.New("resume v1 upload: ger file status error:" + err.Error())
		return
	}

	up := storage.NewResumeUploader(r.cfg)
	err = up.Put(workspace.GetContext(), &ret, info.TokenProvider(), info.SaveKey, file, fileStatus.Size(), r.ext)
	if err != nil {
		err = errors.New("resume v1 upload: upload error:" + err.Error())
	}

	return
}
