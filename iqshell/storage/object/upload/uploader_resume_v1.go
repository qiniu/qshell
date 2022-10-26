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

	if _, sErr := os.Stat(info.FilePath); sErr != nil && os.IsNotExist(sErr) {
		return nil, data.NewEmptyError().AppendDesc("resume v1 upload: get file status error:" + sErr.Error())
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
			return nil, data.NewEmptyError().AppendDesc("resume v1 upload: new recorder error:" + nErr.Error())
		} else {
			recorder = re
		}
	}

	var progress int64 = 0
	ret := &ApiResult{}
	up := storage.NewResumeUploader(r.cfg)
	if pErr := up.PutFile(workspace.GetContext(), &ret, token, info.SaveKey, info.FilePath, &storage.RputExtra{
		Recorder:   recorder,
		Params:     nil,
		UpHost:     info.UpHost,
		MimeType:   info.MimeType,
		TryTimes:   info.TryTimes,
		ChunkSize:  data.BLOCK_SIZE,
		Progresses: nil,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			if info.Progress != nil {
				newProgress := int64(blkIdx) * data.BLOCK_SIZE
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
		return ret, data.NewEmptyError().AppendDesc("resume v1 upload").AppendError(pErr)
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
		return ret, nil
	}
}
