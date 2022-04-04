package api

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"time"
)

type resumeV2 struct {
	ResumeInfo

	uploader *storage.ResumeUploaderV2
}

func (r *resumeV2) InitServer(ctx context.Context) *data.CodeError {
	// uploadId 存在且有效
	if now := time.Now().Unix(); len(r.Recorder.UploadId) > 0 && now < r.Recorder.ExpireTime {
		return nil
	}

	hasKey := len(r.Key) != 0
	ret := &storage.InitPartsRet{}
	err := r.uploader.InitParts(ctx, r.TokenProvider(), r.UpHost, r.Bucket,
		r.Key, hasKey, ret)
	if err == nil {
		r.Recorder.UploadId = ret.UploadID
		r.Recorder.ExpireTime = time.Now().Unix() + 3600*24*5
	} else {
		err = errors.New("resume v2 init server error:" + err.Error())
	}
	return data.ConvertError(err)
}

func (r *resumeV2) UploadBlock(ctx context.Context, index int, d []byte) *data.CodeError {
	hasKey := len(r.Key) != 0
	partNumber := int64(len(r.Recorder.Parts)) + 1
	size := len(d)
	partMd5 := md5.Sum(d)
	partMd5String := hex.EncodeToString(partMd5[:])
	ret := &storage.UploadPartsRet{}
	err := r.uploader.UploadParts(ctx, r.TokenProvider(), r.UpHost, r.Bucket,
		r.Key, hasKey, r.Recorder.UploadId, partNumber, partMd5String, ret, bytes.NewReader(d), size)
	if err == nil {
		r.Recorder.Parts = append(r.Recorder.Parts, storage.UploadPartInfo{
			Etag:       ret.Etag,
			PartNumber: partNumber,
		})
		r.Recorder.Offset += int64(size)
	} else {
		err = errors.New("resume v2 upload block error:" + err.Error())
	}
	return data.ConvertError(err)
}

func (r *resumeV2) Complete(ctx context.Context, putRet interface{}) *data.CodeError {
	hasKey := len(r.Key) != 0
	putExtra := &storage.RputV2Extra{
		Progresses: r.Recorder.Parts,
	}
	err := r.uploader.CompleteParts(ctx, r.TokenProvider(), r.UpHost, &putRet, r.Bucket,
		r.Key, hasKey, r.Recorder.UploadId, putExtra)
	if err != nil {
		err = errors.New("resume v2 complete error:" + err.Error())
	}
	return data.ConvertError(err)
}
