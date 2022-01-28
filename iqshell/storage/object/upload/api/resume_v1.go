package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
)

type resumeV1 struct {
	ResumeInfo

	uploader *storage.ResumeUploader
}

func (r *resumeV1) InitServer(ctx context.Context) error {
	return nil
}

// UploadBlock size 必须是 4M 整数倍
func (r *resumeV1) UploadBlock(ctx context.Context, index int, data []byte) error {
	size := len(data)
	var blkCtx storage.BlkputRet
	err := r.uploader.Mkblk(ctx, r.UpToken, r.UpHost, &blkCtx, size, bytes.NewReader(data), size)
	if err == nil {
		r.Recorder.BlkCtxs = append(r.Recorder.BlkCtxs, blkCtx)
		r.Recorder.Offset += int64(size)
	} else {
		err = errors.New("resume v1 upload block error:" + err.Error())
	}
	return err
}

func (r *resumeV1) Complete(ctx context.Context, putRet interface{}) (err error) {
	putExtra := storage.RputExtra{
		Progresses: r.Recorder.BlkCtxs,
	}
	err = r.uploader.Mkfile(ctx, r.UpToken, r.UpHost, putRet, r.Key, true, r.Recorder.TotalSize, &putExtra)
	if err != nil {
		err = errors.New("resume v1 complete error:" + err.Error())
	}
	return
}
