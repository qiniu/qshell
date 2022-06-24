package ip

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"net/http"
	"strings"
)

var (
	// IP信息查询接口地址
	aliIPParseUrls = []string{
		"https://ip.taobao.com/outGetIpInfo",
		"https://ip.taobao.com/service/getIpInfo.php",
	}
)

type aliParser struct {
}

var _ Parser = (*aliParser)(nil)

func NewAliIPParser() Parser {
	return &aliParser{}
}

func (a *aliParser) Parse(ip string) (result ParserResult, err *data.CodeError) {
	for _, url := range aliIPParseUrls {
		result, err = a.parse(ip, url)
		if err == nil && result != nil {
			break
		}
	}
	return
}

func (a *aliParser) parse(ip string, fromUrl string) (ParserResult, *data.CodeError) {
	req, err := http.NewRequest("GET", fromUrl, nil)
	if err != nil {
		return nil, data.ConvertError(err)
	}

	q := req.URL.Query()
	q.Add("accessKey", "alibaba-inc")
	q.Add("ip", ip)
	req.URL.RawQuery = q.Encode()

	gResp, gErr := http.DefaultClient.Do(req)
	if gErr != nil {
		return nil, data.NewEmptyError().AppendDescF("Query ip info failed for %s, %s", ip, gErr)
	}
	defer gResp.Body.Close()
	responseBody, rErr := io.ReadAll(gResp.Body)
	if rErr != nil {
		return nil, data.NewEmptyError().AppendDescF("read body failed for %s, %s", ip, rErr)
	}

	log.Debug("a parser")

	info := &aliIpInfo{}
	decodeErr := json.Unmarshal(responseBody, info)
	if decodeErr != nil {
		return nil, data.NewEmptyError().AppendDescF("Parse ip failed for %s", ip)
	}

	return info, nil
}

type aliIpInfo struct {
	Code int       `json:"code"`
	Data aliIpData `json:"data"`
}

func (i *aliIpInfo) String() string {
	return fmt.Sprintf("%v", i.Data.String())
}

// IpData ip 具体的信息
type aliIpData struct {
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

func (i aliIpData) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "IP", i.Ip))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Country", i.Country))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Area", i.Area))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Region", i.Region))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "City", i.City))
	s.WriteString(fmt.Sprintf("%-10s:%s", "Isp", i.Isp))
	return s.String()
}
