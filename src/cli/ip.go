package cli

import (
	"fmt"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/log"
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
		client := rs.NewEx(nil)
		for _, ip := range params {
			url := fmt.Sprintf("%s?ip=%s", TAOBAO_IP_QUERY, ip)
			var ipInfo IpInfo
			err := client.Conn.Call(nil, &ipInfo, url)
			if err != nil {
				log.Error("Query ip info failed for", ip, "due to", err)
			} else {
				fmt.Println(fmt.Sprintf("Ip: %-20s => %s", ip, ipInfo))
			}
			<-time.After(time.Second * 1)
		}
	} else {
		CmdHelp(cmd)
	}
}
