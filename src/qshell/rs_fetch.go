package qshell

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

func Fetch(mac *digest.Mac, remoteResUrl, bucket, key string) (err error) {
	client := rs.New(mac)
	fetchUri := fmt.Sprintf("/fetch/%s/to/%s",
		base64.URLEncoding.EncodeToString([]byte(remoteResUrl)),
		base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)))
	err = client.Conn.Call(nil, nil, conf.IO_HOST+fetchUri)
	return
}

func Prefetch(mac *digest.Mac, bucket, key string) (err error) {
	client := rs.New(mac)
	prefetchUri := fmt.Sprintf("/prefetch/%s", base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)))
	err = client.Conn.Call(nil, nil, conf.IO_HOST+prefetchUri)
	return
}
