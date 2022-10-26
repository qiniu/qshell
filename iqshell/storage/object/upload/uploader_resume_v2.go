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

	if _, sErr := os.Stat(info.FilePath); sErr != nil && os.IsNotExist(sErr) {
		return nil, data.NewEmptyError().AppendDesc("resume v2 upload: get file status error:" + sErr.Error())
	}

	token := info.TokenProvider()
	log.DebugF("upload token:%s", token)

	if info.Progress != nil {
		info.Progress.SetFileSize(info.LocalFileSize)
		info.Progress.Start()
	}

	var recorder storage.Recorder = nil
	if len(info.CacheDir) > 0 {
		if re, nErr := storage.NewFileRecorder(info.CacheDir); nErr != nil {
			return nil, data.NewEmptyError().AppendDesc("resume v2 upload: new recorder error:" + nErr.Error())
		} else {
			recorder = re
		}
	}

	var progress int64 = 0
	ret := &ApiResult{}
	up := storage.NewResumeUploaderV2(r.cfg)
	if pErr := up.PutFile(workspace.GetContext(), &ret, token, info.SaveKey, info.FilePath, &storage.RputV2Extra{
		Recorder:   recorder,
		Metadata:   nil,
		CustomVars: nil,
		UpHost:     info.UpHost,
		MimeType:   info.MimeType,
		PartSize:   info.ChunkSize,
		TryTimes:   info.TryTimes,
		Progresses: nil,
		Notify: func(partNumber int64, ret *storage.UploadPartsRet) {
			if info.Progress != nil {
				newProgress := partNumber * info.ChunkSize
				if progress == 0 {
					progress = newProgress
				} else if newProgress-progress >= info.ChunkSize {
					progress += info.ChunkSize
				}
				info.Progress.Progress(progress)
			}
		},
		NotifyErr: nil,
	}); pErr != nil {
		return ret, data.NewEmptyError().AppendDesc("resume v2 upload").AppendError(pErr)
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
		return ret, nil
	}
}
