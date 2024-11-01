package bucket

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storagev2/apis"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func GetStorageV2() (storageClient *apis.Storage, err *data.CodeError) {
	acc, gErr := workspace.GetAccount()
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("GetStorageV2: get current account error:%v", gErr)
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	options := workspace.GetHttpClientOptions()
	options.Credentials = mac
	storageClient = apis.NewStorage(options)
	return
}

func GetHttpClient() (httpClient *http_client.Client, err *data.CodeError) {
	acc, gErr := workspace.GetAccount()
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("GetStorageV2: get current account error:%v", gErr)
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	options := workspace.GetHttpClientOptions()
	options.Credentials = mac
	httpClient = http_client.NewClient(options)
	return
}
