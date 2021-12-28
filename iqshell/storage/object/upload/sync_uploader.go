package upload

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/qiniu/go-sdk/v7/storage"
	storage2 "github.com/qiniu/qshell/v2/iqshell/storage"
	"time"
)

type ResumeUploader struct {
	*storage.ResumeUploader
}

func (p *ResumeUploader) UpHost(ak, bucket string) (upHost string, err error) {
	return storage2.GetUpHost(p.Cfg, ak, bucket)
}

// NewResumeUploader 表示构建一个新的分片上传的对象
func NewResumeUploader(cfg *storage.Config) *ResumeUploader {
	rUploader := storage.NewResumeUploader(cfg)
	return &ResumeUploader{
		ResumeUploader: rUploader,
	}
}

type IResumeUploader interface {
	initServer(ctx context.Context) error
	uploadBlock(ctx context.Context, data []byte) error
	complete(ctx context.Context) (putRet SputRet, err error)
}

type resumeUploaderV1 struct {
	uploader *storage.ResumeUploader
	recorder *ProgressRecorder
	uptoken  string
	key      string
	upHost   string
}

func (uploader *resumeUploaderV1) initServer(ctx context.Context) error {
	return nil
}

// size 必须是 4M 整数倍
func (uploader *resumeUploaderV1) uploadBlock(ctx context.Context, data []byte) error {
	size := len(data)
	var blkCtx storage.BlkputRet
	err := uploader.uploader.Mkblk(ctx, uploader.uptoken, uploader.upHost, &blkCtx, size, bytes.NewReader(data), size)
	if err == nil {
		uploader.recorder.BlkCtxs = append(uploader.recorder.BlkCtxs, blkCtx)
		uploader.recorder.Offset += int64(size)
	}
	return err
}

func (uploader *resumeUploaderV1) complete(ctx context.Context) (putRet SputRet, err error) {
	putExtra := storage.RputExtra{
		Progresses: uploader.recorder.BlkCtxs,
	}
	err = uploader.uploader.Mkfile(ctx, uploader.uptoken, uploader.upHost, &putRet, uploader.key, true, uploader.recorder.TotalSize, &putExtra)
	return
}

type resumeUploaderV2 struct {
	uploader *storage.ResumeUploaderV2
	recorder *ProgressRecorder
	upHost   string
	bucket   string
	uptoken  string
	key      string
}

func (uploader *resumeUploaderV2) initServer(ctx context.Context) error {
	// uploadId 存在且有效
	if now := time.Now().Unix(); len(uploader.recorder.UploadId) > 0 && now < uploader.recorder.ExpireTime {
		return nil
	}

	hasKey := len(uploader.key) != 0
	ret := &storage.InitPartsRet{}
	err := uploader.uploader.InitParts(ctx, uploader.uptoken, uploader.upHost, uploader.bucket,
		uploader.key, hasKey, ret)
	if err == nil {
		uploader.recorder.UploadId = ret.UploadID
		uploader.recorder.ExpireTime = time.Now().Unix() + 3600*24*5
	}
	return err
}

func (uploader *resumeUploaderV2) uploadBlock(ctx context.Context, data []byte) error {
	hasKey := len(uploader.key) != 0
	partNumber := int64(len(uploader.recorder.Parts)) + 1
	size := len(data)
	partMd5 := md5.Sum(data)
	partMd5String := hex.EncodeToString(partMd5[:])
	ret := &storage.UploadPartsRet{}
	err := uploader.uploader.UploadParts(ctx, uploader.uptoken, uploader.upHost, uploader.bucket,
		uploader.key, hasKey, uploader.recorder.UploadId, partNumber, partMd5String, ret, bytes.NewReader(data), size)
	if err == nil {
		uploader.recorder.Parts = append(uploader.recorder.Parts, storage.UploadPartInfo{
			Etag:       ret.Etag,
			PartNumber: partNumber,
		})
		uploader.recorder.Offset += int64(size)
	}
	return err
}

func (uploader *resumeUploaderV2) complete(ctx context.Context) (putRet SputRet, err error) {
	hasKey := len(uploader.key) != 0
	putExtra := &storage.RputV2Extra{
		Progresses: uploader.recorder.Parts,
	}
	err = uploader.uploader.CompleteParts(ctx, uploader.uptoken, uploader.upHost, &putRet, uploader.bucket,
		uploader.key, hasKey, uploader.recorder.UploadId, putExtra)
	return
}
