package iqshell

import (
	"github.com/qiniu/api.v7/v7/storage"
)

type ResumeUploader struct {
	*storage.ResumeUploader
}

func (p *ResumeUploader) UpHost(ak, bucket string) (upHost string, err error) {
	return GetUpHost(p.Cfg, ak, bucket)
}

// NewResumeUploader 表示构建一个新的分片上传的对象
func NewResumeUploader(cfg *storage.Config) *ResumeUploader {
	rUploader := storage.NewResumeUploader(cfg)
	return &ResumeUploader{
		ResumeUploader: rUploader,
	}
}
