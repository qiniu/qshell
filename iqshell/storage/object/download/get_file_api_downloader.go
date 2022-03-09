package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"net/http"
	"strings"
)

type getFileApiDownloader struct {
	useHttps bool
	mac      *qbox.Mac
}

func (g *getFileApiDownloader) Download(info ApiInfo) (response *http.Response, err error) {
	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}
	return g.download(info)
}

func (g *getFileApiDownloader) download(info ApiInfo) (response *http.Response, err error) {
	if len(info.Domain) == 0 {
		info.Domain = workspace.GetConfig().Hosts.GetOneIo()
	}
	if len(info.Domain) == 0 {
		zone, err := bucket.Region(info.Bucket)
		if err != nil {
			return nil, err
		}
		info.Domain = zone.IovipHost
	}

	if len(info.Domain) == 0 {
		if len(info.Host) == 0 {
			err = errors.New("get file api: can't get io domain")
			return
		}
		info.Domain = info.Host
	}

	// /getfile/<ak>/<bucket>/<UrlEncodedKey>[?e=<Deadline>&token=<DownloadToken>
	url := utils.Endpoint(g.useHttps, info.Domain)
	url = strings.Join([]string{url, "getfile", g.mac.AccessKey, info.Bucket, info.Key}, "/")
	url, err = PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
		PublicUrl: url,
		Deadline:  7 * 24 * 3600,
	})
	if err != nil {
		return nil, fmt.Errorf("PublicUrlToPrivate error:%v", err)
	}

	log.DebugF("get file api download, url:%s", url)
	//new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("New request failed by url:" + url + " error:" + err.Error())
	}

	// 设置断点续传
	if info.FromBytes > 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", info.FromBytes))
	}
	// 配置 referer
	if len(info.Referer) > 0 {
		req.Header.Add("Referer", info.Referer)
	}
	log.DebugF("request:\n%+v", req)
	return http.DefaultClient.Do(req)
}
