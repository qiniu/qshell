package operations

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"time"
)

type UploadInfo struct {
	FilePath string
	Bucket   string
	Key      string
	MimeType string
}

func UploadFile(info UploadInfo) {
	doneSignal := make(chan bool)
	go func(ch chan bool) {
		progressSigns := []string{"|", "/", "-", "\\", "|"}
		for {
			for _, p := range progressSigns {
				log.Info("\rProgress: ", p)
				os.Stdout.Sync()
				select {
				case <-ch:
					return
				case <-time.After(time.Millisecond * 50):
					continue
				}
			}
		}
	}(doneSignal)
	ret, err := uploadFile(upload.ApiInfo{
		FilePath: info.FilePath,
		ToBucket: info.Bucket,
		SaveKey:  info.Key,
		MimeType: info.MimeType,
	})
	doneSignal <- true

	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			log.ErrorF("Upload file error %d: %s, Reqid: %s", v.Code, v.Err, v.Reqid)
		}
	} else {
		log.Alert("FileHash:", ret.Hash)
		log.Alert("Fsize:", ret.FSize, "(", utils.FormatFileSize(ret.FSize), ")")
		log.Alert("MimeType:", ret.MimeType)
	}

	if err != nil {
		os.Exit(data.STATUS_ERROR)
	}
}

func uploadFile(info upload.ApiInfo) (res upload.ApiResult, err error) {
	startTime := time.Now().UnixNano() / 1e6
	cfg := workspace.GetConfig()
	uploadConfig := cfg.Up
	info.PutThreshold = uploadConfig.PutThreshold
	info.UseResumeV2 = uploadConfig.IsResumableAPIV2()
	info.ChunkSize = uploadConfig.ResumableAPIV2PartSize
	info.UpHost = uploadConfig.UpHost
	info.DisableResume = uploadConfig.IsDisableResume()
	info.DisableForm = uploadConfig.IsDisableForm()
	if info.TokenProvider == nil {
		info.TokenProvider, err = createTokenProvider(info)
	}
	if err != nil {
		log.ErrorF("Upload  failed because get token provider error:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
		return
	}

	res, err = upload.Upload(info)
	if err != nil {
		log.ErrorF("Upload  failed:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
		return
	}
	endTime := time.Now().UnixNano() / 1e6

	duration := float64(endTime-startTime) / 1000
	speed := fmt.Sprintf("%.2fKB/s", float64(res.FSize)/duration/1024)
	if res.IsSkip {
		log.AlertF("Upload skip because file exist:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	} else {
		log.AlertF("Upload File success %s => [%s:%s] duration:%.2fs speed:%s", info.FilePath, info.ToBucket, info.SaveKey, duration, speed)
	}

	return res, nil
}

func createTokenProvider(info upload.ApiInfo) (provider func() string, err error) {
	mac, gErr := workspace.GetMac()
	if gErr != nil {
		return nil, errors.New("get mac error:" + gErr.Error())
	}

	provider = createTokenProviderWithMac(mac, *workspace.GetConfig().Up.Policy, info)
	return
}

func createTokenProviderWithMac(mac *qbox.Mac, policy storage.PutPolicy, info upload.ApiInfo) func() string {
	policy.Scope = info.ToBucket
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
