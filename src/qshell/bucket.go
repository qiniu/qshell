package qshell

import (
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
)

func GetBuckets(mac *digest.Mac) (buckets []string, err error) {
	buckets = make([]string, 0)
	client := rs.New(mac)
	bucketsUri := fmt.Sprintf("%s/buckets", conf.RS_HOST)
	err = client.Conn.Call(nil, &buckets, bucketsUri)
	return
}

func GetDomainsOfBucket(mac *digest.Mac, bucket string) (domains []string, err error) {
	domains = make([]string, 0)
	client := rs.New(mac)
	getDomainsUrl := fmt.Sprintf("%s/v6/domain/list", DEFAULT_API_HOST)
	postData := map[string][]string{
		"tbl": []string{bucket},
	}
	err = client.Conn.CallWithForm(nil, &domains, getDomainsUrl, postData)
	return
}
