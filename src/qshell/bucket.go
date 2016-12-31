package qshell

import (
	"fmt"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	"qiniu/api.v6/rs"
)

type BucketInfo struct {
	Region string `json:"region"`
}

func GetBucketInfo(mac *digest.Mac, bucket string) (bucketInfo BucketInfo, err error) {
	client := rs.NewMac(mac)
	bucketUri := fmt.Sprintf("%s/bucket/%s", conf.RS_HOST, bucket)
	err = client.Conn.Call(nil, &bucketInfo, bucketUri)
	return
}

func GetBuckets(mac *digest.Mac) (buckets []string, err error) {
	buckets = make([]string, 0)
	client := rs.NewMac(mac)
	bucketsUri := fmt.Sprintf("%s/buckets", conf.RS_HOST)
	err = client.Conn.Call(nil, &buckets, bucketsUri)
	return
}

func GetDomainsOfBucket(mac *digest.Mac, bucket string) (domains []string, err error) {
	domains = make([]string, 0)
	client := rs.NewMac(mac)
	getDomainsUrl := fmt.Sprintf("%s/v6/domain/list", conf.API_HOST)
	postData := map[string][]string{
		"tbl": []string{bucket},
	}
	err = client.Conn.CallWithForm(nil, &domains, getDomainsUrl, postData)
	return
}
