package download

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/host"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
	"net/url"
	"strings"
)

type getFileApiDownloader struct {
	useHttps bool
	mac      *qbox.Mac
}

func (g *getFileApiDownloader) Download(info *ApiInfo) (response *http.Response, err *data.CodeError) {
	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}

	h, hErr := info.HostProvider.Provide()
	if hErr != nil {
		return nil, hErr.HeaderInsertDesc("[provide host]")
	}

	response, err = g.download(h, info)
	if err != nil || (response != nil && response.StatusCode/100 != 2) {
		log.DebugF("download freeze host:%s", h.GetServer())
		info.HostProvider.Freeze(h)
	}

	return response, err
}

func (g *getFileApiDownloader) download(host *host.Host, info *ApiInfo) (*http.Response, *data.CodeError) {

	// /getfile/<ak>/<bucket>/<UrlEncodedKey>[?e=<Deadline>&token=<DownloadToken>
	urlString := utils.Endpoint(g.useHttps, host.GetServer())
	urlString = strings.Join([]string{urlString, "getfile", g.mac.AccessKey, info.Bucket, url.PathEscape(info.Key)}, "/")
	result, err := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
		PublicUrl: urlString,
		Deadline:  7 * 24 * 3600,
	})

	if result == nil || err != nil {
		return nil, data.NewEmptyError().AppendDescF("PublicUrlToPrivate error:%v", err)
	}
	urlString = result.Url

	log.DebugF("get file api download, url:%s", urlString)
	log.DebugF("get download, host:%s", host.GetHost())

	headers := http.Header{}
	// set host
	if len(host.GetHost()) > 0 {
		headers.Add("Host", host.GetHost())
	}

	// 设置断点续传
	if info.FromBytes >= 0 && info.ToBytes >= 0 && (info.FromBytes+info.ToBytes) > 0 {
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

	if workspace.IsCmdInterrupt() {
		return nil, data.CancelError
	}

	response, rErr := defaultClient.DoRequest(workspace.GetContext(), "GET", urlString, headers)
	if len(info.ServerFileHash) != 0 && response != nil && response.Header != nil {
		etag := response.Header.Get("Etag")
		etag = utils.ParseEtag(etag)
		if len(etag) > 0 && etag != info.ServerFileHash {
			return nil, data.NewEmptyError().AppendDescF("file has change, hash before:%s now:%s", info.ServerFileHash, etag)
		}
	}

	if rErr == nil {
		return response, nil
	}

	cErr := data.ConvertError(rErr)
	if cErr.Code <= 0 {
		cErr.Code = data.ErrorCodeUnknown
	}
	return response, cErr
}
