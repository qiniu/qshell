package upload

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
)

type resumeV2Uploader struct {
	cfg *storage.Config
}

func newResumeV2Uploader(cfg *storage.Config) Uploader {
	return &resumeV2Uploader{
		cfg: cfg,
	}
}

func (r *resumeV2Uploader) upload(info *ApiInfo) (*ApiResult, *data.CodeError) {
	log.DebugF("resume v2 upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)

	file, err := os.Open(info.FilePath)
	if err != nil {
		return nil, data.NewEmptyError().AppendDesc("resume v2 upload: open file error:" + err.Error())
	}

	fileStatus, err := file.Stat()
	if err != nil {
		return nil, data.NewEmptyError().AppendDesc("resume v2 upload: ger file status error:" + err.Error())
	}

	token := info.TokenProvider()
	log.DebugF("upload token:%s", token)

	if info.Progress != nil {
		info.Progress.SetFileSize(info.FileSize)
		info.Progress.Start()
	}

	ret := &ApiResult{}
	up := storage.NewResumeUploaderV2(r.cfg)
	err = up.Put(workspace.GetContext(), ret, token, info.SaveKey, file, fileStatus.Size(), &storage.RputV2Extra{
		Recorder:   nil,
		Metadata:   nil,
		CustomVars: nil,
		UpHost:     info.UpHost,
		MimeType:   info.MimeType,
		PartSize:   info.ChunkSize,
		TryTimes:   info.TryTimes,
		Progresses: nil,
		Notify: func(partNumber int64, ret *storage.UploadPartsRet) {
			if info.Progress != nil {
				info.Progress.SendSize(info.ChunkSize)
			}
		},
		NotifyErr: nil,
	})
	if err != nil {
		err = data.NewEmptyError().AppendDesc("resume v2 upload: upload error:" + err.Error())
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
	}

	return ret, data.ConvertError(err)
}
