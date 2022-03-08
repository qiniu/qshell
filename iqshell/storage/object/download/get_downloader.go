package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"net/http"
)

type getDownloader struct {
	useHttps bool
}

func (g *getDownloader) Download(info ApiInfo) (response *http.Response, err error) {
	if len(info.Domain) == 0  {
		log.DebugF("download: get domain of bucket:%s", info.Bucket)
		if d, e := bucket.DomainOfBucket(info.Bucket); e != nil {
			err = fmt.Errorf("download: get bucket domain error:%v, domain can't be empty", e)
			return
		} else {
			info.Domain = d
			log.DebugF("download: bucket:%s domain:%s", info.Bucket, info.Domain)
		}
	}

	if len(info.Domain) == 0 {
		err = errors.New("download: get bucket domain: can't get bucket domain")
		return
	}

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
	if len(info.Domain) > 0 {
		req.Host = info.Domain
	}
	// 设置断点续传
	if info.FromBytes > 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", info.FromBytes))
	}
	// 配置 referer
	if len(info.Referer) > 0 {
		req.Header.Add("Referer", info.Referer)
	}
	return http.DefaultClient.Do(req)
}
