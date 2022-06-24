package upload

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
)

type resumeV1Uploader struct {
	cfg *storage.Config
}

func newResumeV1Uploader(cfg *storage.Config) Uploader {
	return &resumeV1Uploader{
		cfg: cfg,
	}
}

func (r *resumeV1Uploader) upload(info *ApiInfo) (*ApiResult, *data.CodeError) {
	log.DebugF("resume v1 upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)

	file, err := os.Open(info.FilePath)
	if err != nil {
		return nil, data.NewEmptyError().AppendDesc("resume v1 upload: open file error:" + err.Error())
	}

	fileStatus, err := file.Stat()
	if err != nil {
		return nil, data.NewEmptyError().AppendDesc("resume v1 upload: ger file status error:" + err.Error())
	}

	token := info.TokenProvider()
	log.DebugF("upload token:%s", token)

	if info.Progress != nil {
		info.Progress.SetFileSize(info.FileSize)
		info.Progress.Start()
	}

	ret := &ApiResult{}
	up := storage.NewResumeUploader(r.cfg)
	err = up.Put(workspace.GetContext(), &ret, token, info.SaveKey, file, fileStatus.Size(), &storage.RputExtra{
		Recorder:   nil,
		Params:     nil,
		UpHost:     info.UpHost,
		MimeType:   info.MimeType,
		TryTimes:   info.TryTimes,
		Progresses: nil,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			if info.Progress != nil {
				info.Progress.SendSize(int64(blkSize))
			}
		},
		NotifyErr: nil,
	})
	if err != nil {
		err = data.NewEmptyError().AppendDesc("resume v1 upload: upload error:" + err.Error())
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
	}

	return ret, data.ConvertError(err)
}
