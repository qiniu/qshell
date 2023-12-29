package upload

import (
	"os"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qiniu/qshell/v2/iqshell/common/client"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
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
	c := client.DefaultStorageClient()
	up := storage.NewResumeUploaderEx(r.cfg, &c)
	extra := &storage.RputExtra{
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
	}

	var pErr error
	if info.SequentialReadFile {
		file, oErr := os.Open(info.FilePath)
		if oErr != nil {
			return nil, data.NewEmptyError().AppendDesc("resume v1 upload: open error:" + oErr.Error())
		}
		defer file.Close()

		log.Debug("resume v1 upload: put with reader")
		pErr = up.PutWithoutSize(workspace.GetContext(), &ret, token, info.SaveKey, file, extra)
	} else {
		log.Debug("resume v1 upload: put with file path")
		pErr = up.PutFile(workspace.GetContext(), &ret, token, info.SaveKey, info.FilePath, extra)
	}

	if pErr != nil {
		return ret, data.NewEmptyError().AppendDesc("resume v1 upload").AppendError(pErr)
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
		return ret, nil
	}
}
