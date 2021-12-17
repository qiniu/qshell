package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type TokenInfo struct {
	AccessKey         string
	SecretKey         string
	Url               string
	Method            string
	ContentType       string
	Body              string
	PutPolicyFilePath string
}

// CreateQBoxToken QBox Token, 一般bucket相关的接口需要这个token
func CreateQBoxToken(info TokenInfo) {
	mac, req, mErr := getMacAndRequest(info)
	if mErr != nil {
		log.ErrorF("create mac and request: %v\n", mErr)
		os.Exit(data.STATUS_ERROR)
	}
	token, signErr := mac.SignRequest(req)
	if signErr != nil {
		log.ErrorF("create qbox token for url: %s: %v\n", info.Url, signErr)
		os.Exit(data.STATUS_ERROR)
	}
	log.Alert("QBox " + token)
}

// 签名七牛token, 一般三鉴接口需要http头文件 Authorization, 这个头的值就是qiniuToken
func CreateQiniuToken(info TokenInfo) {
	mac, req, mErr := getMacAndRequest(info)
	if mErr != nil {
		log.ErrorF("create mac and reqeust: %v\n", mErr)
		os.Exit(data.STATUS_ERROR)
	}
	token, signErr := mac.SignRequestV2(req)
	if signErr != nil {
		log.ErrorF("create qiniu token for url: %s: %v\n", info.Url, signErr)
		os.Exit(data.STATUS_ERROR)
	}
	log.Alert("Qiniu " + token)
}

// 给定上传策略，打印出上传token
func CreateUploadToken(info TokenInfo) {
	fileName := info.PutPolicyFilePath

	fileInfo, oErr := os.Open(fileName)
	if oErr != nil {
		log.ErrorF("open file %s: %v\n", fileName, oErr)
		os.Exit(data.STATUS_ERROR)
	}
	defer fileInfo.Close()

	configData, rErr := ioutil.ReadAll(fileInfo)
	if rErr != nil {
		log.ErrorF("read putPolicy config file `%s`: %v\n", fileName, rErr)
		os.Exit(data.STATUS_ERROR)
	}

	//remove UTF-8 BOM
	configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))

	putPolicy := new(storage.PutPolicy)
	uErr := json.Unmarshal(configData, putPolicy)
	if uErr != nil {
		log.ErrorF("parse upload config file `%s`: %v\n", fileName, uErr)
		os.Exit(data.STATUS_ERROR)
	}

	var mac *qbox.Mac
	var mErr error
	if info.AccessKey != "" && info.SecretKey != "" {
		mac = qbox.NewMac(info.AccessKey, info.SecretKey)
	} else {
		mac, mErr = account.GetMac()
		if mErr != nil {
			log.ErrorF("get mac: %v\n", mErr)
			os.Exit(1)
		}
	}
	uploadToken := putPolicy.UploadToken(mac)
	fmt.Println("UpToken " + uploadToken)
}

func getMacAndRequest(info TokenInfo) (mac *qbox.Mac, req *http.Request, err error) {

	var mErr error
	var rErr error

	if info.AccessKey != "" && info.SecretKey != "" {
		mac = qbox.NewMac(info.AccessKey, info.SecretKey)
	} else {
		mac, mErr = account.GetMac()
		if mErr != nil {
			err = fmt.Errorf("get mac: %v\n", mErr)
			return
		}
	}
	var httpBody io.Reader
	if info.Body == "" {
		httpBody = nil
	} else {
		httpBody = strings.NewReader(info.Body)
	}
	req, rErr = http.NewRequest(info.Method, info.Url, httpBody)
	if rErr != nil {
		err = fmt.Errorf("create request: %v\n", rErr)
		return
	}
	headers := http.Header{}
	headers.Add("Content-Type", info.ContentType)
	req.Header = headers

	return
}
