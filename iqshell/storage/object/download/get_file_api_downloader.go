package download

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
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
	return g.download(info)
}

func (g *getFileApiDownloader) download(info *ApiInfo) (*http.Response, *data.CodeError) {
	host, hErr := info.HostProvider.Provide()
	if hErr != nil {
		return nil, hErr.HeaderInsertDesc("[provide host]")
	}

	// /getfile/<ak>/<bucket>/<UrlEncodedKey>[?e=<Deadline>&token=<DownloadToken>
	url := utils.Endpoint(g.useHttps, host.GetServer())
	url = strings.Join([]string{url, "getfile", g.mac.AccessKey, info.Bucket, info.Key}, "/")
	url, err := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
		PublicUrl: url,
		Deadline:  7 * 24 * 3600,
	})
	if err != nil {
		return nil, data.NewEmptyError().AppendDescF("PublicUrlToPrivate error:%v", err)
	}

	log.DebugF("get file api download, url:%s", url)
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
