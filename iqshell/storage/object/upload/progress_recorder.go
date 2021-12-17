package upload

import "github.com/qiniu/go-sdk/v7/storage"

type ProgressRecorder struct {
	BlkCtxs      []storage.BlkputRet      `json:"blk_ctxs"`    // resume v1
	Parts        []storage.UploadPartInfo `json:"parts"`       // resume v2
	UploadId     string                   `json:"upload_id"`   // resume v2
	ExpireTime   int64                    `json:"expire_time"` // resume v2
	Offset       int64                    `json:"offset"`
	TotalSize    int64                    `json:"total_size"`
	LastModified int                      `json:"last_modified"` // 上传文件的modification time
	FilePath     string                   // 断点续传记录保存文件
}

func NewProgressRecorder(filePath string) *ProgressRecorder {
	p := new(ProgressRecorder)
	p.FilePath = filePath
	p.BlkCtxs = make([]storage.BlkputRet, 0)
	p.Parts = make([]storage.UploadPartInfo, 0)
	return p
}