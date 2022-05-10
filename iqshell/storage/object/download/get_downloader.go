package download

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
)

type getDownloader struct {
	useHttps bool
}

func (g *getDownloader) Download(info *ApiInfo) (*http.Response, *data.CodeError) {
	host, hErr := info.HostProvider.Provide()
	if hErr != nil {
		return nil, hErr.HeaderInsertDesc("[provide host]")
	}

	url := ""
	// 构造下载 url
	if info.IsPublic {
		url = PublicUrl(UrlApiInfo{
			BucketDomain: host.GetServer(),
			Key:          info.Key,
			UseHttps:     g.useHttps,
		})
	} else {
		url = PrivateUrl(UrlApiInfo{
			BucketDomain: host.GetServer(),
			Key:          info.Key,
			UseHttps:     g.useHttps,
		})
	}

	log.DebugF("get download, url:%s", url)
	log.DebugF("get download, host:%s", host.GetHost())

	headers := http.Header{}
	// set host
	if len(host.GetHost()) > 0 {
		headers.Add("Host", host.GetHost())
	}

	// 设置断点续传
	if info.FromBytes > 0 {
		headers.Add("Range", fmt.Sprintf("bytes=%d-", info.FromBytes))
	}
	// 配置 referer
	if len(info.Referer) > 0 {
		headers.Add("Referer", info.Referer)
	}
	response, rErr := storage.DefaultClient.DoRequest(workspace.GetContext(), "GET", url, headers)
	if utils.IsHostUnavailableError(rErr) {
		log.DebugF("download freeze host:%s because: %v", host.GetServer(), rErr)
		info.HostProvider.Freeze(host)
	}
	return response, data.ConvertError(rErr)
}
