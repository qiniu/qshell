package qshell

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qiniu/api.v6/rs"
)

func BatchRefresh(client *rs.Client, urls []string) (err error) {
	if len(urls) == 0 || len(urls) > 100 {
		err = errors.New("url count invalid, should between [1, 100]")
		return
	}

	postUrl := "http://fusion.qiniuapi.com/refresh"

	postData := map[string][]string{
		"urls": urls,
	}

	err = client.Conn.CallWithForm(nil, nil, postUrl, postData)
	return
}

type CdnIpInfo struct {
	LineCname string `json:"LineCname"`
	CdnInfo   string `json:"cdninfo"`
	IpAddress string `json:"ipaddress"`
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
