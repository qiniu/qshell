package qshell

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qiniu/api.v6/rs"
)

type BatchRefreshRequest struct {
	Urls []string `json:"urls"`
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

func BatchRefresh(client *rs.Client, urls []string) (batchRefreshResp BatchRefreshResponse, err error) {
	if len(urls) == 0 || len(urls) > 100 {
		err = errors.New("url count invalid, should between [1, 100]")
		return
	}

	postUrl := "http://fusion.qiniuapi.com/v2/tune/refresh"

	batchRefreshReq := BatchRefreshRequest{
		Urls: urls,
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

type CdnIpInfo struct {
	LineCname []string `json:"LineCname"`
	CdnInfo   string   `json:"cdninfo"`
	IpAddress string   `json:"ipaddress"`
}

func GetCdnSupplierOfIp(client *rs.Client, ip string) (cdnInfo CdnIpInfo, err error) {
	getUrl := fmt.Sprintf("http://api.qiniu.com/v1/ipcdninfo/%s", ip)
	var resp *http.Response
	resp, err = client.Conn.Get(nil, getUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, rErr := ioutil.ReadAll(resp.Body)
	if rErr != nil {
		err = rErr
		return
	}

	decodeErr := json.Unmarshal(respBody, &cdnInfo)
	if decodeErr != nil {
		err = decodeErr
		return
	}

	return
}

func GetCdnRegionalIps(client *rs.Client, cname, isp, province string) (ips []string, err error) {
	postUrl := "http://api.qiniu.com/v1/regionalip/"
	postData := map[string][]string{
		"cname":    []string{cname},
		"isp":      []string{isp},
		"province": []string{province},
	}
	ips = make([]string, 0, 200)
	err = client.Conn.CallWithForm(nil, &ips, postUrl, postData)
	return
}
