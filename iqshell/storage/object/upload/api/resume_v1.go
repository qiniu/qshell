package api

import (
	"bytes"
	"context"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type resumeV1 struct {
	ResumeInfo

	uploader *storage.ResumeUploader
}

func (r *resumeV1) InitServer(ctx context.Context) *data.CodeError {
	return nil
}

// UploadBlock size 必须是 4M 整数倍
func (r *resumeV1) UploadBlock(ctx context.Context, index int, d []byte) *data.CodeError {
	size := len(d)
	var blkCtx storage.BlkputRet
	err := r.uploader.Mkblk(ctx, r.TokenProvider(), r.UpHost, &blkCtx, size, bytes.NewReader(d), size)
	if err != nil {
		return data.NewEmptyError().AppendDesc("resume v1 upload block error:" + err.Error())
	} else {
		r.Recorder.BlkCtxs = append(r.Recorder.BlkCtxs, blkCtx)
		r.Recorder.Offset += int64(size)
		return nil
	}
}

func (r *resumeV1) Complete(ctx context.Context, putRet interface{}) *data.CodeError {
	putExtra := storage.RputExtra{
		Progresses: r.Recorder.BlkCtxs,
	}
	if err := r.uploader.Mkfile(ctx, r.TokenProvider(), r.UpHost, putRet, r.Key, true, r.Recorder.TotalSize, &putExtra); err != nil {
		return data.NewEmptyError().AppendDescF("resume v1 complete error:%v", err)
	}
	return nil
}
