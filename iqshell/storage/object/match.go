package object

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

const (
	MatchCheckModeFileHash = 0
	MatchCheckModeFileSize = 1
)

type MatchApiInfo struct {
	Bucket         string `json:"bucket"`     // 文件所在七牛云的空间名，【必选】
	Key            string `json:"key"`        // 文件所在七牛云的 Key， 【必选】
	LocalFile      string `json:"local_file"` // 本地文件路径；【必选】
	CheckMode      int    `json:"-"`          // 检测模式， 0: 检测 hash，其他检测 size 【可选】
	ServerFileHash string `json:"file_hash"`  // 服务端文件 Etag，可以是 etagV1, 也可以是 etagV2；【可选】没有会从服务获取
	ServerFileSize int64  `json:"file_size"`  // 服务端文件大小；【可选】没有会从服务获取
}

func (m *MatchApiInfo) WorkId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", m.Bucket, m.Key, m.LocalFile, m.ServerFileHash))
}

func (m *MatchApiInfo) CheckModeHash() bool {
	return m.CheckMode == MatchCheckModeFileHash
}

type MatchResult struct {
	Exist bool
	Match bool
}

var _ flow.Result = (*MatchResult)(nil)

func (m *MatchResult) IsValid() bool {
	return m != nil && m.Match
}

func Match(info MatchApiInfo) (match *MatchResult, err *data.CodeError) {
	if len(info.LocalFile) == 0 {
		return &MatchResult{
			Exist: false,
			Match: false,
		}, data.NewEmptyError().AppendDesc("Match Check, file is empty")
	}

	if info.CheckModeHash() {
		return matchHash(info)
	} else {
		return matchSize(info)
	}
}

func matchSize(info MatchApiInfo) (result *MatchResult, err *data.CodeError) {
	result = &MatchResult{
		Exist: false,
		Match: false,
	}

	if info.ServerFileSize <= 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: false,
		}); sErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match check size, get file status").AppendError(sErr)
		} else {
			info.ServerFileSize = stat.FSize
			result.Exist = true
		}
	} else {
		result.Exist = true
	}

	stat, sErr := os.Stat(info.LocalFile)
	if sErr != nil {
		return result, data.NewEmptyError().AppendDescF("Match check size, get local file status").AppendError(sErr)
	}
	if info.ServerFileSize == stat.Size() {
		result.Match = true
		return result, nil
	} else {
		result.Match = false
		return result, data.NewEmptyError().AppendDescF("Match check size, size don't match, file:%s except:%d but:%d", info.LocalFile, stat.Size(), info.ServerFileSize)
	}
}

func matchHash(info MatchApiInfo) (result *MatchResult, err *data.CodeError) {
	result = &MatchResult{
		Exist: false,
		Match: false,
	}

	hashFile, oErr := os.Open(info.LocalFile)
	if oErr != nil {
		return result, data.NewEmptyError().AppendDescF("Match check hash, get local file error:%v", oErr)
	}

	var serverObjectStat *StatusResult
	if len(info.ServerFileHash) == 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: true,
		}); sErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match Check, get file status").AppendError(sErr)
		} else {
			info.ServerFileHash = stat.Hash
			serverObjectStat = &stat
			result.Exist = true
		}
	} else {
		result.Exist = true
	}

	// 计算本地文件 hash
	var hash string
	if utils.IsSignByEtagV2(info.ServerFileHash) {
		log.DebugF("Match check hash: get etag by v2 for key:%s", info.Key)
		if serverObjectStat == nil {
			if stat, sErr := Status(StatusApiInfo{
				Bucket:   info.Bucket,
				Key:      info.Key,
				NeedPart: true,
			}); sErr != nil {
				return result, data.NewEmptyError().AppendDescF("Match check hash, etag v2, get file status").AppendError(sErr)
			} else {
				serverObjectStat = &stat
			}
		}
		if h, eErr := utils.EtagV2(hashFile, serverObjectStat.Parts); eErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match check hash, get file etag v2").AppendError(eErr)
		} else {
			hash = h
		}
		log.DebugF("Match check hash, get etag by v2 for key:%s hash:%s", info.Key, hash)
	} else {
		log.DebugF("Match check hash, get etag by v1 for key:%s", info.Key)
		if h, eErr := utils.EtagV1(hashFile); eErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match check hash, get file etag v1").AppendError(eErr)
		} else {
			hash = h
		}
		log.DebugF("Match check hash, get etag by v1 for key:%s hash:%s", info.Key, hash)
	}
	log.DebugF("Match check hash,       server hash, key:%s hash:%s", info.Key, hash)
	if hash != info.ServerFileHash {
		return result, data.NewEmptyError().AppendDescF("Match check hash, file hash doesn't match for key:%s, local file hash:%s server file hash:%s", info.Key, hash, info.ServerFileHash)
	}

	result.Match = true
	return result, nil
}
