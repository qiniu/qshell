package rs

import (
	"encoding/base64"
	"net/http"

	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
)

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

// @gist entry
type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Customer string `json:"customer"`
}

// @endgist

func (rs Client) Stat(l rpc.Logger, bucket, key string) (entry Entry, err error) {
	err = rs.Conn.Call(l, &entry, RS_HOST+URIStat(bucket, key))
	return
}

func (rs Client) Delete(l rpc.Logger, bucket, key string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIDelete(bucket, key))
}

func (rs Client) Move(l rpc.Logger, bucketSrc, keySrc, bucketDest, keyDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIMove(bucketSrc, keySrc, bucketDest, keyDest))
}

func (rs Client) Copy(l rpc.Logger, bucketSrc, keySrc, bucketDest, keyDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URICopy(bucketSrc, keySrc, bucketDest, keyDest))
}

func (rs Client) ChangeMime(l rpc.Logger, bucket, key, mime string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIChangeMime(bucket, key, mime))
}

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func URIDelete(bucket, key string) string {
	return "/delete/" + encodeURI(bucket+":"+key)
}

func URIStat(bucket, key string) string {
	return "/stat/" + encodeURI(bucket+":"+key)
}

func URICopy(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/copy/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIMove(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/move/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIChangeMime(bucket, key, mime string) string {
	return "/chgm/" + encodeURI(bucket+":"+key) + "/mime/" + encodeURI(mime)
}
