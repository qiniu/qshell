package cli

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
	"time"
)

const (
	TAOBAO_IP_QUERY = "http://ip.taobao.com/service/getIpInfo.php"
)

type IpInfo struct {
	Code int    `json:"code"`
	Data IpData `json:"data"`
}

func (this IpInfo) String() string {
	return fmt.Sprintf("%s", this.Data)
}

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

func IpQuery(cmd string, params ...string) {
	if len(params) > 0 {
		for _, ip := range params {
			url := fmt.Sprintf("%s?ip=%s", TAOBAO_IP_QUERY, ip)
			var ipInfo IpInfo
			func() {
				gResp, gErr := http.Get(url)
				if gErr != nil {
					logs.Error("Query ip info failed for %s, %s", ip, gErr)
					return
				}
				defer gResp.Body.Close()
				//fmt.Println(fmt.Sprintf("Ip: %-20s => %s", ip, ipInfo))
				decoder := json.NewDecoder(gResp.Body)
				decodeErr := decoder.Decode(&ipInfo)
				if decodeErr != nil {
					logs.Error("Parse ip info body failed for %s, %s", ip, decodeErr)
					return
				}

				fmt.Println(fmt.Sprintf("%s\t%s", ip, ipInfo.String()))
			}()
			<-time.After(time.Millisecond * 500)
		}
	} else {
		CmdHelp(cmd)
	}
}
