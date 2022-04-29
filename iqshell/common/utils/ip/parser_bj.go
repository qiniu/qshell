package ip

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const (
	// IP信息查询接口地址
	bjIPParseUrl = "https://www.bejson.com/Bejson/Api/Ip/getIp"
)

type bjParser struct {
}

var _ Parser = (*bjParser)(nil)

func NewBjIPParser() Parser {
	return &bjParser{}
}

func (a bjParser) Parse(ip string) (ParserResult, *data.CodeError) {
	req, err := http.NewRequest("GET", bjIPParseUrl, nil)
	if err != nil {
		return nil, data.ConvertError(err)
	}

	q := req.URL.Query()
	q.Add("ip", ip)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("User-Agent", getUA())

	gResp, gErr := http.DefaultClient.Do(req)
	if gErr != nil {
		return nil, data.NewEmptyError().AppendDescF("Query ip info failed for %s, %s", ip, gErr)
	}
	defer gResp.Body.Close()
	responseBody, rErr := io.ReadAll(gResp.Body)
	if rErr != nil {
		return nil, data.NewEmptyError().AppendDescF("read body failed for %s, %s", ip, rErr)
	}
	log.Debug("b parser")

	info := &bjIpInfo{}
	decodeErr := json.Unmarshal(responseBody, info)
	if decodeErr != nil {
		return nil, data.NewEmptyError().AppendDescF("Parse ip failed for %s", ip)
	}

	return info, nil
}

type bjIpInfo struct {
	Code int      `json:"code"`
	Data bjIpData `json:"data"`
}

func (i *bjIpInfo) String() string {
	return fmt.Sprintf("%v", i.Data)
}

// IpData ip 具体的信息
type bjIpData struct {
	Country string `json:"country"`
	//CountryId string `json:"country_id"`
	Area string `json:"area"`
	//AreaId    string `json:"area_id"`
	Region string `json:"region"`
	//RegionId  string `json:"region_id"`
	City string `json:"city"`
	//CityId    string `json:"city_id"`
	Isp string `json:"isp"`
	//IspId     int64  `json:"isp_id"`
	Ip string `json:"ip"`
}

func (i bjIpData) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "IP", i.Ip))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Country", i.Country))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Area", i.Area))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "Region", i.Region))
	s.WriteString(fmt.Sprintf("%-10s:%s\n", "City", i.City))
	s.WriteString(fmt.Sprintf("%-10s:%s", "Isp", i.Isp))
	return s.String()
}

func getUA() string {
	m1 := rand.Intn(2) + 10
	m2 := rand.Intn(2)
	if m1 < 11 {
		m2 = rand.Intn(16)
	}
	os := fmt.Sprintf("Macintosh; Intel Mac OS X %d_%d_%d", m1, m2, rand.Intn(2))
	c := fmt.Sprintf("Chrome/%d.%d.4280.%d", rand.Intn(10)+75, rand.Intn(1), rand.Intn(100))
	if rand.Int()%2 == 0 {
		c += fmt.Sprintf(" Safari/%d.%d", rand.Intn(50)+480, rand.Intn(50))
	}
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) %s", os, c)
}
