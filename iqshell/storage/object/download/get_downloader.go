package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"net/http"
)

type getDownloader struct {
	useHttps bool
}

func (g *getDownloader) Download(info *ApiInfo) (response *http.Response, err error) {
	url := ""
	// 构造下载 url
	if info.IsPublic {
		url = PublicUrl(UrlApiInfo{
			BucketDomain: info.Domain,
			Key:          info.Key,
			UseHttps:     g.useHttps,
		})
	} else {
		url = PrivateUrl(UrlApiInfo{
			BucketDomain: info.Domain,
			Key:          info.Key,
			UseHttps:     g.useHttps,
		})
	}

	log.DebugF("get download, url:%s", url)
	//new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("New request failed by url:" + url + " error:" + err.Error())
	}
	// set host
	if len(info.Host) > 0 {
		req.Host = info.Host
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
