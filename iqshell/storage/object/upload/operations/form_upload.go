package operations

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"strings"
	"time"
)

type FormUploadInfo struct {
	FilePath     string
	Host         string
	Bucket       string
	Key          string
	FileType     int
	Overwrite    bool
	MimeType     string
	CallbackUrls string
	CallbackHost string
}

type formPutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

// FormUpload 【fput】使用表单上传本地文件到七牛存储空间
func FormUpload(info FormUploadInfo) {
	if info.FileType != 1 && info.FileType != 0 {
		log.Error("Wrong Filetype, It should be 0 or 1")
		os.Exit(data.STATUS_ERROR)
	}

	//create uptoken
	policy := storage.PutPolicy{}
	if info.Overwrite {
		policy.Scope = fmt.Sprintf("%s:%s", info.Bucket, info.Key)
	} else {
		policy.Scope = info.Bucket
	}
	policy.FileType = info.FileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	if (info.CallbackUrls == "" && info.CallbackHost != "") || (info.CallbackUrls != "" && info.CallbackHost == "") {
		log.Error("callbackUrls and callback must exist at the same time")
		os.Exit(1)
	}
	if info.CallbackHost != "" && info.CallbackUrls != "" {
		info.CallbackUrls = strings.Replace(info.CallbackUrls, ",", ";", -1)
		policy.CallbackHost = info.CallbackHost
		policy.CallbackURL = info.CallbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}

	var putExtra storage.PutExtra
	var upHost string

	if info.Host == "" {
		if len(workspace.GetConfig().Hosts.Up) > 0 {
			upHost = workspace.GetConfig().Hosts.Up[0]
		}
	} else {
		upHost = info.Host
	}
	putExtra = storage.PutExtra{
		UpHost: upHost,
	}
	if info.MimeType != "" {
		putExtra.MimeType = info.MimeType
	}

	mac, err := account.GetMac()
	if err != nil {
		log.Error("Get Mac error: ", err)
		os.Exit(data.STATUS_ERROR)
	}

	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := formPutRet{}
	startTime := time.Now()
	fStat, statErr := os.Stat(info.FilePath)
	if statErr != nil {
		log.ErrorF("Local file error: %v", statErr)
		os.Exit(data.STATUS_ERROR)
	}
	fsize := fStat.Size()
	log.ErrorF("Uploading %s => %s : %s ...\n", info.FilePath, info.Bucket, info.Key)

	formUploader := storage.NewFormUploader(nil)

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

	err = formUploader.PutFile(workspace.GetContext(), &putRet, uptoken, info.Key, info.FilePath, &putExtra)

	doneSignal <- true

	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			log.ErrorF("Put file error %d: %s, Reqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			log.ErrorF("Put file error: %v", err)
		}
	} else {
		log.Info("\rProgress: 100%")
		log.Alert("")

		log.Alert("Put file", info.FilePath, "=>", info.Bucket, ":", putRet.Key, "success!")
		log.Alert("FileHash:", putRet.Hash)
		log.Alert("Fsize:", putRet.Fsize, "(", utils.FormatFileSize(fsize), ")")
		log.Alert("MimeType:", putRet.MimeType)

		lastNano := time.Now().UnixNano() - startTime.UnixNano()
		lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
		avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
		log.Alert("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
	}

	if err != nil {
		os.Exit(data.STATUS_ERROR)
	}
}
