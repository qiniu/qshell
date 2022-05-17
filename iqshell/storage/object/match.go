package object

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

type MatchApiInfo struct {
	Bucket    string // 文件所在七牛云的空间名，必选
	Key       string // 文件所在七牛云的 Key， 必选
	FileHash  string // 文件 Etag，可以是 etagV1, 也可以是 etagV2；可选，没有会从服务获取
	LocalFile string // 本地文件路径；必选
}

func (m *MatchApiInfo) WorkId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", m.Bucket, m.Key, m.LocalFile, m.FileHash))
}

type MatchResult struct {
	Match bool
}

var _ flow.Result = (*MatchResult)(nil)

func (m *MatchResult) IsValid() bool {
	return m != nil
}

func Match(info MatchApiInfo) (match *MatchResult, err *data.CodeError) {
	if len(info.LocalFile) == 0 {
		return nil, data.NewEmptyError().AppendDesc("Match Check: file is empty")
	}

	hashFile, oErr := os.Open(info.LocalFile)
	if oErr != nil {
		return nil, data.NewEmptyError().AppendDescF("Match Check: get local file error:%v", oErr)
	}

	var serverObjectStat *StatusResult
	if len(info.FileHash) == 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: true,
		}); sErr != nil {
			return nil, data.NewEmptyError().AppendDescF("Match Check: get file status error:%v", sErr)
		} else {
			info.FileHash = stat.Hash
			serverObjectStat = &stat
		}
	}

	// 计算本地文件 hash
	var hash string
	if utils.IsSignByEtagV2(info.FileHash) {
		log.DebugF("Match Check: get etag by v2 for key:%s", info.Key)
		if serverObjectStat == nil {
			if stat, sErr := Status(StatusApiInfo{
				Bucket:   info.Bucket,
				Key:      info.Key,
				NeedPart: true,
			}); sErr != nil {
				return nil, data.NewEmptyError().AppendDescF("Match Check: etag v2, get file status error:%v", sErr)
			} else {
				serverObjectStat = &stat
			}
		}
		if h, eErr := utils.EtagV2(hashFile, serverObjectStat.Parts); eErr != nil {
			return nil, data.NewEmptyError().AppendDescF("Match Check: get file etag v2 error:%v", eErr)
		} else {
			hash = h
		}
		log.DebugF("Match Check: get etag by v2 for key:%s hash:%s", info.Key, hash)
	} else {
		log.DebugF("Match Check: get etag by v1 for key:%s", info.Key)
		if h, eErr := utils.EtagV1(hashFile); eErr != nil {
			return nil, data.NewEmptyError().AppendDescF("Match Check: get file etag v1 error:%v", eErr)
		} else {
			hash = h
		}
		log.DebugF("Match Check: get etag by v1 for key:%s hash:%s", info.Key, hash)
	}
	log.DebugF("Match Check:       server hash, key:%s hash:%s", info.Key, hash)
	if hash != info.FileHash {
		return nil, data.NewEmptyError().AppendDescF("Match Check: file hash doesn't match for key:%s, local file hash:%s server file hash:%s", info.Key, hash, info.FileHash)
	}

	return &MatchResult{
		Match: true,
	}, nil
}
