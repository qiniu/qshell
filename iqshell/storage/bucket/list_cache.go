package bucket

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

type cacheInfo struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
	Marker string `json:"marker"`
}

type listCache struct {
	enableRecord bool
	cachePath    string
}

func (l *listCache) saveCache(info *cacheInfo) *data.CodeError {
	if !l.enableRecord || info == nil {
		return nil
	}

	if len(l.cachePath) == 0 {
		return data.NewError(0, "load cache: no cache path set, will not save record")
	}

	return utils.MarshalToFile(l.cachePath, info)
}

func (l *listCache) loadCache() (info *cacheInfo, err *data.CodeError) {
	if !l.enableRecord {
		return nil, nil
	}

	if len(l.cachePath) == 0 {
		return nil, data.NewError(0, "load cache: no cache path set, will not load record")
	}

	info = &cacheInfo{}
	err = utils.UnMarshalFromFile(l.cachePath, info)
	return
}

func (l *listCache) removeCache() *data.CodeError {
	if !l.enableRecord || len(l.cachePath) == 0 {
		return nil
	}

	err := os.Remove(l.cachePath)
	return data.ConvertError(err)
}
