package qshell

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

func GetBuckets(mac *digest.Mac) (buckets []string, err error) {
	buckets = make([]string, 0)
	client := rs.New(mac)
	bucketsUri := fmt.Sprintf("%s/buckets", conf.RS_HOST)
	err = client.Conn.Call(nil, &buckets, bucketsUri)
	return
}
