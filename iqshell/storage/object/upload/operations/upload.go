package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"io/ioutil"
	"os"
	"strings"
)

type BatchUploadInfo struct {
	FailExportFilePath     string
	SuccessExportFilePath  string
	OverrideExportFilePath string
	UpThreadCount          int64
	ConfigFile             string
	UploadConfig           config.UploadConfig
}

// [qupload]命令， 上传本地文件到七牛存储中

// BatchUpload 该命令会读取配置文件， 上传本地文件系统的文件到七牛存储中;
// 可以设置多线程上传，默认的线程区间在[iqshell.MIN_UPLOAD_THREAD_COUNT, iqshell.MAX_UPLOAD_THREAD_COUNT]
func BatchUpload(info BatchUploadInfo) {

	if len(info.ConfigFile) > 0 {
		pErr := parseUploadConfigFile(info.ConfigFile, &(info.UploadConfig))
		if pErr != nil {
			log.Error(fmt.Sprintf("Upload parse config file: %s: %v\n", info.ConfigFile, pErr))
			os.Exit(data.STATUS_HALT)
		}
	}

	uploadConfig := info.UploadConfig
	if uploadConfig.FileType != 1 && uploadConfig.FileType != 0 {
		log.Error("Wrong Filetype, It should be 0 or 1 ")
		os.Exit(data.STATUS_HALT)
	}

	srcFileInfo, err := os.Stat(uploadConfig.SrcDir)
	if err != nil {
		log.Error("Upload config error for parameter `SrcDir`,", err)
		os.Exit(data.STATUS_HALT)
	}

	if !srcFileInfo.IsDir() {
		log.Error("Upload src dir should be a directory")
		os.Exit(data.STATUS_HALT)
	}

	if uploadConfig.Bucket == "" {
		fmt.Println("Upload config no `bucket` specified")
		os.Exit(data.STATUS_HALT)
	}

	policy := storage.PutPolicy{}

	if (uploadConfig.CallbackUrls == "" && uploadConfig.CallbackHost != "") || (uploadConfig.CallbackUrls != "" && uploadConfig.CallbackHost == "") {
		log.Error("callbackUrls and callback must exist at the same time")
		os.Exit(1)
	}

	if uploadConfig.CallbackHost != "" && uploadConfig.CallbackUrls != "" {
		callbackUrls := strings.Replace(uploadConfig.CallbackUrls, ",", ";", -1)
		policy.CallbackHost = uploadConfig.CallbackHost
		policy.CallbackURL = callbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}
	uploadConfig.PutPolicy = policy

	//upload
	if info.UpThreadCount < upload.MIN_UPLOAD_THREAD_COUNT || info.UpThreadCount > upload.MAX_UPLOAD_THREAD_COUNT {
		log.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			upload.MIN_UPLOAD_THREAD_COUNT, upload.MAX_UPLOAD_THREAD_COUNT)

		if info.UpThreadCount < upload.MIN_UPLOAD_THREAD_COUNT {
			info.UpThreadCount = upload.MIN_UPLOAD_THREAD_COUNT
		} else if info.UpThreadCount > upload.MAX_UPLOAD_THREAD_COUNT {
			info.UpThreadCount = upload.MAX_UPLOAD_THREAD_COUNT
		}
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.SuccessExportFilePath,
		FailExportFilePath:     info.FailExportFilePath,
		OverrideExportFilePath: info.OverrideExportFilePath,
	})
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}
	upload.QiniuUpload(int(info.UpThreadCount), &uploadConfig, resultExport)
}

func parseUploadConfigFile(uploadConfigFile string, uploadConfig *config.UploadConfig) (err error) {
	//read upload config
	if uploadConfigFile == "" {
		err = fmt.Errorf("config filename is empty")
		return
	}
	fp, oErr := os.Open(uploadConfigFile)
	if oErr != nil {
		err = fmt.Errorf("Open upload config file ``%s`: %v\n", uploadConfigFile, oErr)
		return
	}
	defer fp.Close()

	configData, rErr := ioutil.ReadAll(fp)
	if rErr != nil {
		err = fmt.Errorf("Read upload config file `%s`: %v\n", uploadConfigFile, rErr)
		return
	}

	//remove UTF-8 BOM
	configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))
	uErr := json.Unmarshal(configData, uploadConfig)
	if uErr != nil {
		err = fmt.Errorf("Parse upload config file `%s`: %v\n", uploadConfigFile, uErr)
		return
	}
	return
}
