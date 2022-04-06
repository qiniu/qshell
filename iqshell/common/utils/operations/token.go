package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type TokenInfo struct {
	AccessKey   string
	SecretKey   string
	Url         string
	Method      string
	Body        string
	ContentType string
}

type QBoxTokenInfo struct {
	TokenInfo
}

func (info *QBoxTokenInfo) Check() *data.CodeError {
	if len(info.Url) == 0 {
		return alert.CannotEmptyError("Url", "")
	}
	return nil
}

// Token QBox Token, 一般bucket相关的接口需要这个token
func Token(cfg *iqshell.Config) {
	iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{})
}

// CreateQBoxToken QBox Token, 一般bucket相关的接口需要这个token
func CreateQBoxToken(cfg *iqshell.Config, info QBoxTokenInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	mac, req, mErr := getMacAndRequest(info.TokenInfo)
	if mErr != nil {
		log.ErrorF("create mac and request: %v\n", mErr)
		os.Exit(data.StatusError)
	}
	token, signErr := mac.SignRequest(req)
	if signErr != nil {
		log.ErrorF("create qbox token for url: %s: %v\n", info.Url, signErr)
		os.Exit(data.StatusError)
	}
	log.Alert("QBox " + token)
}

type QiniuTokenInfo struct {
	TokenInfo
}

func (info *QiniuTokenInfo) Check() *data.CodeError {
	if len(info.Url) == 0 {
		return alert.CannotEmptyError("Url", "")
	}
	return nil
}

// CreateQiniuToken 签名七牛token, 一般三鉴接口需要http头文件 Authorization, 这个头的值就是qiniuToken
func CreateQiniuToken(cfg *iqshell.Config, info QiniuTokenInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	mac, req, mErr := getMacAndRequest(info.TokenInfo)
	if mErr != nil {
		log.ErrorF("create mac and reqeust: %v\n", mErr)
		os.Exit(data.StatusError)
	}
	token, signErr := mac.SignRequestV2(req)
	if signErr != nil {
		log.ErrorF("create qiniu token for url: %s: %v\n", info.Url, signErr)
		os.Exit(data.StatusError)
	}
	log.Alert("Qiniu " + token)
}

type UploadTokenInfo struct {
	TokenInfo
	PutPolicyFilePath string
}

func (info *UploadTokenInfo) Check() *data.CodeError {
	if len(info.PutPolicyFilePath) == 0 {
		return alert.CannotEmptyError("PutPolicyConfigFile", "")
	}
	return nil
}

// 给定上传策略，打印出上传token
func CreateUploadToken(cfg *iqshell.Config, info UploadTokenInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	fileName := info.PutPolicyFilePath

	fileInfo, oErr := os.Open(fileName)
	if oErr != nil {
		log.ErrorF("open file %s: %v\n", fileName, oErr)
		os.Exit(data.StatusError)
	}
	defer fileInfo.Close()

	configData, rErr := ioutil.ReadAll(fileInfo)
	if rErr != nil {
		log.ErrorF("read putPolicy config file `%s`: %v\n", fileName, rErr)
		os.Exit(data.StatusError)
	}

	//remove UTF-8 BOM
	configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))

	putPolicy := new(storage.PutPolicy)
	uErr := json.Unmarshal(configData, putPolicy)
	if uErr != nil {
		log.ErrorF("parse upload config file `%s`: %v\n", fileName, uErr)
		os.Exit(data.StatusError)
	}

	var mac *qbox.Mac
	var mErr *data.CodeError
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

func getMacAndRequest(info TokenInfo) (mac *qbox.Mac, req *http.Request, err *data.CodeError) {
	if info.AccessKey != "" && info.SecretKey != "" {
		mac = qbox.NewMac(info.AccessKey, info.SecretKey)
	} else {
		var mErr error
		mac, mErr = account.GetMac()
		if mErr != nil {
			err = data.NewEmptyError().AppendDescF("get mac: %v\n", mErr)
			return
		}
	}

	var httpBody io.Reader
	if info.Body == "" {
		httpBody = nil
	} else {
		httpBody = strings.NewReader(info.Body)
	}

	var rErr error
	req, rErr = http.NewRequest(info.Method, info.Url, httpBody)
	if rErr != nil {
		err = data.NewEmptyError().AppendDescF("create request: %v\n", rErr)
		return
	}
	headers := http.Header{}
	headers.Add("Content-Type", info.ContentType)
	req.Header = headers

	return
}
