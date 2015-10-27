package rsf

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
)

// ----------------------------------------------------------

type ListItem struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}

type ListRet struct {
	Marker string     `json:"marker"`
	Items  []ListItem `json:"items"`
}

// ----------------------------------------------------------

type Client struct {
	Conn rpc.Client
}

func New(mac *digest.Mac) Client {
	t := digest.NewTransport(mac, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
}

func NewEx(t http.RoundTripper) Client {
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
}

// ----------------------------------------------------------
// 1. 首次请求 marker = ""
// 2. 无论 err 值如何，均应该先看 entries 是否有内容
// 3. 如果后续没有更多数据，err 返回 EOF，markerOut 返回 ""（但不通过该特征来判断是否结束）
func (rsf Client) ListPrefix(l rpc.Logger, bucket, prefix, marker string, limit int) (entries []ListItem, markerOut string, err error) {

	if bucket == "" {
		err = errors.New("bucket could not be nil")
		return
	}

	URL := makeListURL(bucket, prefix, marker, limit)
	listRet := ListRet{}
	err = rsf.Conn.Call(l, &listRet, URL)

	if err != nil {
		return
	}
	if listRet.Marker == "" {
		return listRet.Items, "", io.EOF
	}
	return listRet.Items, listRet.Marker, err
}

func makeListURL(bucket, prefix, marker string, limit int) string {

	query := make(url.Values)
	query.Add("bucket", bucket)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}

	return RSF_HOST + "/list?" + query.Encode()
}
