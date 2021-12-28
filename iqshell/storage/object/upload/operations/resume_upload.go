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

type ResumeUploadInfo struct {
	WorkerCount      int
	IsResumeV2       bool
	ResumeV2PartSize int64
	FilePath         string
	Host             string
	Bucket           string
	Key              string
	FileType         int
	Overwrite        bool
	MimeType         string
	CallbackUrls     string
	CallbackHost     string
}

type resumePutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

// ResumeUpload 使用分片上传本地文件到七牛存储空间, 一般用于较大文件的上传
// 文件会被分割成4M大小的块， 一块一块地上传文件
func ResumeUpload(info ResumeUploadInfo) {
	fStat, statErr := os.Stat(info.FilePath)
	if statErr != nil {
		log.ErrorF("Local file error", statErr)
		os.Exit(data.STATUS_ERROR)
	}
	fsize := fStat.Size()

	//create uptoken
	policy := storage.PutPolicy{}
	policy.Scope = fmt.Sprintf("%s:%s", info.Bucket, info.Key)

	if !info.Overwrite {
		policy.InsertOnly = 1
	}
	policy.FileType = info.FileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	if (info.CallbackUrls == "" && info.CallbackHost != "") || (info.CallbackUrls != "" && info.CallbackHost == "") {
		log.Error("callbackUrls and callback must exist at the same time")
		os.Exit(1)
	}
	if info.CallbackHost != "" && info.CallbackUrls != "" {
		callbackUrls := strings.Replace(info.CallbackUrls, ",", ";", -1)
		policy.CallbackHost = info.CallbackHost
		policy.CallbackURL = callbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}

	var upHost string
	if info.Host == "" {
		if len(workspace.GetConfig().Hosts.Up) > 0 {
			upHost = workspace.GetConfig().Hosts.Up[0]
		}
	} else {
		upHost = info.Host
	}

	mac, err := account.GetMac()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get Mac error: ", err)
		os.Exit(data.STATUS_ERROR)
	}
	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := resumePutRet{}
	startTime := time.Now()
	fmt.Printf("Uploading %s => %s : %s ...\n", info.FilePath, info.Bucket, info.Key)

	if info.IsResumeV2 {
		resume_uploader := storage.NewResumeUploaderV2(nil)
		putExtra := storage.RputV2Extra{
			UpHost:   upHost,
			PartSize: info.ResumeV2PartSize,
		}
		if info.MimeType != "" {
			putExtra.MimeType = info.MimeType
		}
		err = resume_uploader.PutFile(workspace.GetContext(), &putRet, uptoken, info.Key, info.FilePath, &putExtra)
	} else {
		resume_uploader := storage.NewResumeUploader(nil)
		putExtra := storage.RputExtra{
			UpHost: upHost,
		}
		if info.MimeType != "" {
			putExtra.MimeType = info.MimeType
		}
		err = resume_uploader.PutFile(workspace.GetContext(), &putRet, uptoken, info.Key, info.FilePath, &putExtra)
	}

	log.Alert("")
	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			log.ErrorF("Put file error %d: %s, Reqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			log.ErrorF("Put file error: %v", err)
		}
	} else {
		log.Alert("Put file", info.FilePath, "=>", info.Bucket, ":", putRet.Key, "success!")
		log.Alert("Hash:", putRet.Hash)
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
