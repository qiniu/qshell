package operations

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"time"
)

type UploadInfo struct {
	upload.ApiInfo
	RelativePathToSrcPath string // 相对与上传文件夹的路径信息
	Policy                *storage.PutPolicy
	DeleteOnSuccess       bool
}

func (info *UploadInfo) Check() *data.CodeError {
	if len(info.ToBucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.SaveKey) == 0 && len(info.FilePath) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.FilePath) == 0 {
		return alert.CannotEmptyError("LocalFile", "")
	}
	if utils.IsNetworkSource(info.FilePath) {
		return alert.Error("file can't be network source", "")
	}
	return nil
}

func (info *UploadInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", info.FilePath, info.ToBucket, info.SaveKey)
}

func UploadFile(cfg *iqshell.Config, info UploadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	info.Progress = progress.NewPrintProgress(" 进度")
	ret, err := uploadFile(&info)
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Upload file error: %v", err)
	} else {
		log.Alert("")
		log.Alert("-------------- File FlowInfo --------------")
		log.AlertF("%10s%s", "Key: ", ret.Key)
		log.AlertF("%10s%s", "Hash: ", ret.ServerFileHash)
		log.AlertF("%10s%d%s", "FileSize: ", ret.ServerFileSize, "("+utils.FormatFileSize(ret.ServerFileSize)+")")
		log.AlertF("%10s%d", "PutTime: ", ret.ServerPutTime)
		log.AlertF("%10s%s", "MimeType: ", ret.MimeType)
	}
}

func uploadFile(info *UploadInfo) (res *upload.ApiResult, err *data.CodeError) {
	startTime := time.Now().UnixNano() / 1e6
	if info.TokenProvider == nil {
		info.TokenProvider, err = createTokenProvider(info)
	}
	if err != nil {
		log.ErrorF("Upload  failed because get token provider error:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
		return
	}

	res, err = upload.Upload(&info.ApiInfo)
	if err != nil {
		log.ErrorF("Upload  failed:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
		return
	}
	endTime := time.Now().UnixNano() / 1e6

	duration := float64(endTime-startTime) / 1000
	speed := fmt.Sprintf("%.2fKB/s", float64(res.ServerFileSize)/duration/1024)
	if res.IsSkip {
		log.AlertF("Upload skip because file exist:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	} else {
		log.AlertF("Upload File success %s => [%s:%s] duration:%.2fs Speed:%s", info.FilePath, info.ToBucket, info.SaveKey, duration, speed)

		//delete on success
		if info.DeleteOnSuccess {
			deleteErr := os.Remove(info.FilePath)
			if deleteErr != nil {
				log.ErrorF("Delete `%s` on upload success error due to `%s`", info.FilePath, deleteErr)
			} else {
				log.InfoF("Delete `%s` on upload success done", info.FilePath)
			}
		}
	}

	return res, nil
}

func createTokenProvider(info *UploadInfo) (provider func() string, err *data.CodeError) {
	mac, gErr := workspace.GetMac()
	if gErr != nil {
		return nil, data.NewEmptyError().AppendDesc("get mac error:" + gErr.Error())
	}

	provider = createTokenProviderWithMac(mac, info)
	return
}

func createTokenProviderWithMac(mac *qbox.Mac, info *UploadInfo) func() string {
	policy := *info.Policy
	policy.Scope = info.ToBucket
	policy.InsertOnly = 1 // 仅新增不覆盖
	if info.Overwrite {
		policy.Scope = fmt.Sprintf("%s:%s", info.ToBucket, info.SaveKey)
		policy.InsertOnly = 0
	}
	policy.ReturnBody = upload.ApiResultFormat()
	policy.FileType = info.FileType
	return func() string {
		policy.Expires = 7 * 24 * 3600
		return policy.UploadToken(mac)
	}
}
