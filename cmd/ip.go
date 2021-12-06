package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/cobra"
)

const (
	// IP信息查询接口地址
	TAOBAO_IP_QUERY = "http://ip.taobao.com/service/getIpInfo.php"
)

// 接口返回的IP信息
type IpInfo struct {
	Code int    `json:"code"`
	Data IpData `json:"data"`
}

func (this IpInfo) String() string {
	return fmt.Sprintf("%s", this.Data)
}

// ip 具体的信息
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

var ipQueryCmd = &cobra.Command{
	Use:   "ip <Ip1> [<Ip2> [<Ip3> ...]]]",
	Short: "Query the ip information",
	Args:  cobra.MinimumNArgs(1),
	Run:   IpQuery,
}

func init() {
	RootCmd.AddCommand(ipQueryCmd)
}

// 【ip】查询ip的相关信息
func IpQuery(cmd *cobra.Command, params []string) {
	for _, ip := range params {

		var ipInfo IpInfo
		func() {
			req, err := http.NewRequest("GET", TAOBAO_IP_QUERY, nil)
			if err != nil {
				log.Error("%v", err)
				return
			}

			q := req.URL.Query()
			q.Add("accessKey", "alibaba-inc")
			q.Add("ip", ip)
			req.URL.RawQuery = q.Encode()

			gResp, gErr := http.DefaultClient.Do(req)
			if gErr != nil {
				log.Error("Query ip info failed for %s, %s", ip, gErr)
				return
			}
			defer gResp.Body.Close()
			//fmt.Println(fmt.Sprintf("Ip: %-20s => %s", ip, ipInfo))
			decoder := json.NewDecoder(gResp.Body)
			decodeErr := decoder.Decode(&ipInfo)
			if decodeErr != nil {
				log.Error("Parse ip info body failed for %s, %s", ip, decodeErr)
				return
			}

			fmt.Println(fmt.Sprintf("%s\t%s", ip, ipInfo.String()))
		}()
		<-time.After(time.Millisecond * 500)
	}
}
