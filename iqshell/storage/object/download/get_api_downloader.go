package download

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

func (g *getApiDownloader) Download(info *ApiInfo) (*http.Response, *data.CodeError) {
	if len(info.Bucket) == 0 {
		return nil, alert.CannotEmptyError("bucket", "")
	}

	if len(info.Key) == 0 {
		return nil, alert.CannotEmptyError("key", "")
	}

	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}
	return g.download(info)
}

func (g *getApiDownloader) download(info *ApiInfo) (*http.Response, *data.CodeError) {
	entryUri := strings.Join([]string{info.Bucket, info.Key}, ":")
	url := strings.Join([]string{info.Domain, "get", utils.Encode(entryUri)}, "/")

	var d struct {
		URL string `json:"url"`
	}
	ctx := auth.WithCredentials(workspace.GetContext(), g.mac)
	headers := http.Header{}

	log.DebugF("get api download, get url:%s", url)
	if e := storage.DefaultClient.Call(ctx, &d, "GET", url, headers); e != nil {
		return nil, data.ConvertError(e)
	}

	log.DebugF("get api download, url:%s", d.URL)
	response, e := storage.DefaultClient.DoRequest(workspace.GetContext(), "GET", d.URL, headers)
	return response, data.ConvertError(e)
}
