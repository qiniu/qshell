package rs

import (
	"encoding/base64"
	"net/http"

	"fmt"
	"qiniu/api.v6/auth/digest"
	. "qiniu/api.v6/conf"
	"qiniu/rpc"
)

// ----------------------------------------------------------

type Client struct {
	Conn rpc.Client
}

func NewMac(mac *digest.Mac) Client {
	t := digest.NewTransport(mac, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client, ""}}
}

func NewEx(t http.RoundTripper) Client {
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client, ""}}
}

func NewMacEx(mac *digest.Mac, t http.RoundTripper, bindRemoteIp string) Client {
	mt := digest.NewTransport(mac, t)
	client := &http.Client{Transport: mt}
	return Client{rpc.Client{client, bindRemoteIp}}
}

// ----------------------------------------------------------

// @gist entry
type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Customer string `json:"customer"`
	FileType int    `json:"type"`
}

// @endgist

func (rs Client) Stat(l rpc.Logger, bucket, key string) (entry Entry, err error) {
	err = rs.Conn.Call(l, &entry, RS_HOST+URIStat(bucket, key))
	return
}

func (rs Client) Delete(l rpc.Logger, bucket, key string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIDelete(bucket, key))
}

func (rs Client) Move(l rpc.Logger, bucketSrc, keySrc, bucketDest, keyDest string, force bool) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIMove(bucketSrc, keySrc, bucketDest, keyDest, force))
}

func (rs Client) Copy(l rpc.Logger, bucketSrc, keySrc, bucketDest, keyDest string, force bool) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URICopy(bucketSrc, keySrc, bucketDest, keyDest, force))
}

func (rs Client) ChangeMime(l rpc.Logger, bucket, key, mime string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+URIChangeMime(bucket, key, mime))
}

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func URIDelete(bucket, key string) string {
	return fmt.Sprintf("/delete/%s", encodeURI(bucket+":"+key))
}

func URIStat(bucket, key string) string {
	return fmt.Sprintf("/stat/%s", encodeURI(bucket+":"+key))
}

func URICopy(bucketSrc, keySrc, bucketDest, keyDest string, force bool) string {
	return fmt.Sprintf("/copy/%s/%s/force/%v", encodeURI(bucketSrc+":"+keySrc), encodeURI(bucketDest+":"+keyDest), force)
}

func URIMove(bucketSrc, keySrc, bucketDest, keyDest string, force bool) string {
	return fmt.Sprintf("/move/%s/%s/force/%v", encodeURI(bucketSrc+":"+keySrc), encodeURI(bucketDest+":"+keyDest), force)
}

func URIChangeMime(bucket, key, mime string) string {
	return fmt.Sprintf("/chgm/%s/mime/%s", encodeURI(bucket+":"+key), encodeURI(mime))
}

func URIPrefetch(bucket, key string) string {
	return fmt.Sprintf("/prefetch/%s", encodeURI(bucket+":"+key))
}
