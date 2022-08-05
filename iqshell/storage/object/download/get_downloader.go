package download

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/host"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
)

type getDownloader struct {
	useHttps bool
}

func (g *getDownloader) Download(info *ApiInfo) (response *http.Response, err *data.CodeError) {
	h, hErr := info.HostProvider.Provide()
	if hErr != nil {
		return nil, hErr.HeaderInsertDesc("[provide host]")
	}

	for i := 0; i < 3; i++ {
		response, err = g.download(h, info)
		if utils.IsHostUnavailableError(err) {
			break
		}

		if response != nil {
			if response.StatusCode/100 == 2 && err == nil {
				break
			}

			if (response.StatusCode > 399 && response.StatusCode < 500) ||
				response.StatusCode == 612 || response.StatusCode == 631 {
				break
			}
		}
	}

	if err != nil || (response != nil && response.StatusCode/100 != 2) {
		if response == nil {
			info.HostProvider.Freeze(h)
			log.DebugF("download freeze host:%s because:%v", h.GetServer(), err)
		} else if response.StatusCode > 499 && response.StatusCode < 600 {
			info.HostProvider.Freeze(h)
			log.DebugF("download freeze host:%s because:[%s] %v", h.GetServer(), response.Status, err)
		} else {
			log.DebugF("download not freeze host:%s because:[%s] %v", h.GetServer(), response.Status, err)
		}
	}

	return response, err
}

func (g *getDownloader) download(host *host.Host, info *ApiInfo) (*http.Response, *data.CodeError) {
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
	return response, data.ConvertError(rErr)
}
