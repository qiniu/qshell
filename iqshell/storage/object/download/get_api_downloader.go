package download

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
	"strings"
)

type getApiDownloader struct {
	useHttps   bool
	mac        *qbox.Mac
	httpClient *client.Client
}

func (g *getApiDownloader) Download(info *ApiInfo) (response *http.Response, err error) {
	if len(info.Bucket) == 0 {
		return nil, errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.Key) == 0 {
		return nil, errors.New(alert.CannotEmpty("key", ""))
	}

	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}
	return g.download(info)
}

func (g *getApiDownloader) download(info *ApiInfo) (response *http.Response, err error) {
	entryUri := strings.Join([]string{info.Bucket, info.Key}, ":")
	url := strings.Join([]string{info.Domain, "get", utils.Encode(entryUri)}, "/")

	var data struct {
		URL string `json:"url"`
	}
	ctx := auth.WithCredentials(workspace.GetContext(), g.mac)
	headers := http.Header{}

	log.DebugF("get api download, get url:%s", url)
	if err := storage.DefaultClient.Call(ctx, &data, "GET", url, headers); err != nil {
		return nil, err
	}

	log.DebugF("get api download, url:%s", data.URL)
	return storage.DefaultClient.DoRequest(workspace.GetContext(), "GET", data.URL, headers)
}
