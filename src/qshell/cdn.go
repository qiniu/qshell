package qshell

import (
	"errors"
	"qshell/qiniu/api.v6/rs"
)

type BatchRefreshRequest struct {
	Urls []string `json:"urls,omitempty"`
	Dirs []string `json:"dirs,omitempty"`
}

type BatchRefreshResponse struct {
	Code          int      `json:"code"`
	Error         string   `json:"error"`
	RequestId     string   `json:"requestId"`
	InvalidUrls   []string `json:"invalidUrls"`
	InvalidDirs   []string `json:"invalidDirs"`
	UrlQuotaDay   int      `json:"urlQuotaDay"`
	UrlSurplusDay int      `json:"urlSurplusDay"`
	DirQuotaDay   int      `json:"dirQuotaDay"`
	DirSurplusDay int      `json:"dirSurplusDay"`
}

func BatchRefresh(client *rs.Client, urls []string, dirs []string) (batchRefreshResp BatchRefreshResponse, err error) {
	if len(urls) > 100 {
		err = errors.New("url count invalid, should between [1, 100]")
		return
	}

	if len(dirs) > 10 {
		err = errors.New("dir count invalid, should between [1, 10]")
		return
	}

	if len(urls) == 0 && len(dirs) == 0 {
		err = errors.New("no url or dir to refresh error")
		return
	}

	postUrl := "http://fusion.qiniuapi.com/v2/tune/refresh"
	batchRefreshReq := BatchRefreshRequest{
		Urls: urls,
		Dirs: dirs,
	}

	batchRefreshResp = BatchRefreshResponse{}
	err = client.Conn.CallWithJson(nil, &batchRefreshResp, postUrl, batchRefreshReq)
	return
}

type BatchPrefetchRequest struct {
	Urls []string `json:"urls"`
}

type BatchPrefetchResponse struct {
	Code        int      `json:"code"`
	Error       string   `json:"error"`
	RequestId   string   `json:"requestId"`
	InvalidUrls []string `json:"invalidUrls"`
	QuotaDay    int      `json:"quotaDay"`
	SurplusDay  int      `json:"surplusDay"`
}

func BatchPrefetch(client *rs.Client, urls []string) (batchPrefetchResponse BatchPrefetchResponse, err error) {
	if len(urls) == 0 || len(urls) > 100 {
		err = errors.New("url count invalid, should between [1, 100]")
		return
	}

	postUrl := "http://fusion.qiniuapi.com/v2/tune/prefetch"

	batchPrefetchReq := BatchRefreshRequest{
		Urls: urls,
	}

	batchPrefetchResponse = BatchPrefetchResponse{}
	err = client.Conn.CallWithJson(nil, &batchPrefetchResponse, postUrl, batchPrefetchReq)
	return
}
