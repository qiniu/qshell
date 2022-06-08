package upload

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

func (f *formUploader) upload(info *ApiInfo) (ret *ApiResult, err *data.CodeError) {
	log.DebugF("form upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)

	file, oErr := os.Open(info.FilePath)
	if oErr != nil {
		err = data.NewEmptyError().AppendDesc("form upload: open file error:" + err.Error())
		return
	}

	fileStatus, sErr := file.Stat()
	if sErr != nil {
		err = data.NewEmptyError().AppendDesc("form upload: ger file status error:" + err.Error())
		return
	}

	token := info.TokenProvider()
	log.DebugF("upload token:%s", token)

	if info.Progress != nil {
		info.Progress.SetFileSize(info.LocalFileSize)
		info.Progress.Start()
		f.ext.OnProgress = func(fsize, uploaded int64) {
			info.Progress.Progress(uploaded)
		}
	}

	up := storage.NewFormUploader(f.cfg)
	if e := up.Put(workspace.GetContext(), &ret, token, info.SaveKey, file, fileStatus.Size(), f.ext); e != nil {
		err = data.NewEmptyError().AppendDesc("form upload: upload error:" + e.Error())
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
	}

	return
}
