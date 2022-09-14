package download

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/host"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net"
	"net/http"
	"time"
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
	if info.FromBytes >= 0 && info.ToBytes >= 0 {
		if info.FromBytes > 0 && info.FromBytes == info.ToBytes {
			return &http.Response{
				Status:     "already download",
				StatusCode: 200,
			}, nil
		} else if info.ToBytes == 0 {
			headers.Add("Range", fmt.Sprintf("bytes=%d-", info.FromBytes))
		} else {
			headers.Add("Range", fmt.Sprintf("bytes=%d-%d", info.FromBytes, info.ToBytes))
		}
	}

	// 配置 referer
	if len(info.Referer) > 0 {
		headers.Add("Referer", info.Referer)
	}
	response, rErr := defaultClient.DoRequest(workspace.GetContext(), "GET", url, headers)
	if response != nil && response.Header != nil {
		etag := response.Header.Get("Etag")
		if len(etag) > 0 && etag != fmt.Sprintf("\"%s\"", info.ServerFileHash) {
			return nil, data.NewEmptyError().AppendDescF("file has change, hash before:%s now:%s", info.ServerFileHash, etag)
		}
	}
	return response, data.ConvertError(rErr)
}

var defaultClient = storage.Client{
	Client: &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
}
