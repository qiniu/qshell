package api

import (
	"context"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type Resume interface {
	InitServer(ctx context.Context) *data.CodeError
	UploadBlock(ctx context.Context, index int, data []byte) *data.CodeError
	Complete(ctx context.Context, ret interface{}) (err *data.CodeError)
}

type ResumeInfo struct {
	UpHost        string
	Bucket        string
	TokenProvider func() string // token provider
	Key           string
	Cfg           *storage.Config
	Recorder      *ProgressRecorder
}

func NewResume(info ResumeInfo, isResumeV2 bool) Resume {
	if isResumeV2 {
		return &resumeV2{
			ResumeInfo: info,
			uploader:   storage.NewResumeUploaderV2(info.Cfg),
		}
	} else {
		return &resumeV1{
			ResumeInfo: info,
			uploader:   storage.NewResumeUploader(info.Cfg),
		}
	}
}

type retryResume struct {
	resume        Resume
	retryMax      int
	retryInterval time.Duration
}

func NewRetryResume(r Resume, retryMax int, retryInterval time.Duration) Resume {
	return &retryResume{
		resume:        r,
		retryMax:      retryMax,
		retryInterval: retryInterval,
	}
}

func (r retryResume) InitServer(ctx context.Context) (err *data.CodeError) {
	retryTimes := 0
	for {
		err = r.resume.InitServer(ctx)
		if err == nil {
			break
		}

		if retryTimes >= r.retryMax {
			return
		}

		time.Sleep(r.retryInterval)
		log.DebugF("resume api Retrying %d time for init server for error:%v", retryTimes, err)
		retryTimes++
	}
	return err
}

func (r retryResume) UploadBlock(ctx context.Context, index int, data []byte) (err *data.CodeError) {
	retryTimes := 0
	for {
		err = r.resume.UploadBlock(ctx, index, data)
		if err != nil && retryTimes >= r.retryMax {
			return
		}
		if err == nil {
			break
		}
		time.Sleep(r.retryInterval)
		log.DebugF("resume api Retrying %d time for upload block index:[%d] for error:%v", retryTimes, index, err)
		retryTimes++
	}
	return err
}

func (r retryResume) Complete(ctx context.Context, ret interface{}) (err *data.CodeError) {
	retryTimes := 0
	for {
		err = r.resume.Complete(ctx, &ret)
		if err != nil && retryTimes >= r.retryMax {
			err = data.NewEmptyError().AppendDescF("resume api complete error:%v", err)
			return
		}
		if err == nil {
			break
		}
		time.Sleep(r.retryInterval)
		log.DebugF("resume api Retrying %d time for server to create file for error:%v", retryTimes, err)
		retryTimes++
	}
	return err
}
