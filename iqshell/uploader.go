package qshell

import (
	"github.com/qiniu/api.v7/storage"
)

type ResumeUploader struct {
	*storage.ResumeUploader
	client *storage.Client
	cfg    *storage.Config
}

func (p *ResumeUploader) UpHost(ak, bucket string) (upHost string, err error) {
	return GetUpHost(p.cfg, ak, bucket)
}

// NewResumeUploader 表示构建一个新的分片上传的对象
func NewResumeUploader(cfg *storage.Config) *ResumeUploader {
	rUploader := storage.NewResumeUploader(cfg)
	if cfg == nil {
		cfg = &storage.Config{}
	}

	return &ResumeUploader{
		cfg:            cfg,
		client:         &storage.DefaultClient,
		ResumeUploader: rUploader,
	}
}
