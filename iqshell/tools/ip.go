package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"net/http"
	"time"
)

const (
	// IP信息查询接口地址
	ipQueryUrl = "http://ip.taobao.com/service/getIpInfo.php"
)

type IpQueryInfo struct {
	Ips []string
}

func IpQuery(info IpQueryInfo) {
	if len(info.Ips) == 0 {
		log.Error(errors.New(alert.CannotEmpty("ip", "")))
		return
	}

	for _, ip := range info.Ips {
		var ipInfo IpInfo
		func() {
			req, err := http.NewRequest("GET", ipQueryUrl, nil)
			if err != nil {
				log.Error(err)
				return
			}

			q := req.URL.Query()
			q.Add("accessKey", "alibaba-inc")
			q.Add("ip", ip)
			req.URL.RawQuery = q.Encode()

			gResp, gErr := http.DefaultClient.Do(req)
			if gErr != nil {
				log.ErrorF("Query ip info failed for %s, %s", ip, gErr)
				return
			}
			defer gResp.Body.Close()
			//fmt.Println(fmt.Sprintf("Ip: %-20s => %s", ip, ipInfo))
			decoder := json.NewDecoder(gResp.Body)
			decodeErr := decoder.Decode(&ipInfo)
			if decodeErr != nil {
				log.ErrorF("Parse ip info body failed for %s, %s", ip, decodeErr)
				return
			}

			fmt.Println(fmt.Sprintf("%s\t%s", ip, ipInfo.String()))
		}()
		<-time.After(time.Millisecond * 500)
	}
}

type IpInfo struct {
	Code int    `json:"code"`
	Data IpData `json:"data"`
}

func (this IpInfo) String() string {
	return fmt.Sprintf("%s", this.Data)
}

// IpData ip 具体的信息
type IpData struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	County    string `json:"county"`
	CountyId  string `json:"county_id"`
	Isp       string `json:"isp"`
	IspId     string `json:"isp_id"`
	Ip        string `json:"ip"`
}

func (this IpData) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
		this.Country, this.Area, this.Region, this.City, this.County, this.Isp)
}
