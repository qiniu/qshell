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
	Bucket    string
	Key       string
	FileHash  string
	LocalFile string
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
		return nil, data.NewEmptyError().AppendDesc("Match Check: get local file error:" + oErr.Error())
	}

	var serverObjectStat *StatusResult
	if len(info.FileHash) == 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket: info.Bucket,
			Key:    info.Key,
		}); sErr != nil {
			return nil, data.NewEmptyError().AppendDesc("Match Check: get file status error:" + sErr.Error())
		} else {
			info.FileHash = stat.Hash
			serverObjectStat = &stat
		}
	}

	// 计算本地文件 hash
	var hash string
	if utils.IsSignByEtagV2(info.FileHash) {
		log.Debug("Match Check: get etag by v2 for key:" + info.Key)
		if serverObjectStat == nil {
			if stat, sErr := Status(StatusApiInfo{
				Bucket: info.Bucket,
				Key:    info.Key,
			}); sErr != nil {
				return nil, data.NewEmptyError().AppendDesc("Match Check: etag v2, get file status error:" + sErr.Error())
			} else {
				serverObjectStat = &stat
			}
		}
		if h, eErr := utils.EtagV2(hashFile, serverObjectStat.Parts); eErr != nil {
			return nil, data.NewEmptyError().AppendDesc("Match Check: get file etag v2 error:" + eErr.Error())
		} else {
			hash = h
		}

	} else {
		log.Debug("Match Check: get etag by v1 for key:" + info.Key)
		if h, eErr := utils.EtagV1(hashFile); eErr != nil {
			return nil, data.NewEmptyError().AppendDesc("Match Check: get file etag v1 error:" + eErr.Error())
		} else {
			hash = h
		}
	}

	if hash != info.FileHash {
		return nil, data.NewEmptyError().AppendDesc("Match Check: file hash doesn't match for key:" + info.Key + "download file hash:" + hash + " excepted:" + info.FileHash)
	}

	return &MatchResult{
		Match: true,
	}, nil
}
