package qshell

import (
	"fmt"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
)

type BucketInfo struct {
	Region string `json:"region"`
}

var (
	BUCKET_RS_HOST  = "http://rs.qiniu.com"
	BUCKET_API_HOST = "http://api.qiniu.com"
)

/*
get bucket info

@param mac
@param bucket - bucket name

@return bucketInfo, err
*/
func GetBucketInfo(mac *digest.Mac, bucket string) (bucketInfo BucketInfo, err error) {
	client := rs.NewMac(mac)
	bucketUrl := fmt.Sprintf("%s/bucket/%s", BUCKET_RS_HOST, bucket)
	callErr := client.Conn.Call(nil, &bucketInfo, bucketUrl)
	if callErr != nil {
		if v, ok := callErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("code: %d, %s, xreqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			err = callErr
		}
	}
	return
}

func GetBuckets(mac *digest.Mac) (buckets []string, err error) {
	buckets = make([]string, 0)
	client := rs.NewMac(mac)
	bucketsUri := fmt.Sprintf("%s/buckets", BUCKET_RS_HOST)
	callErr := client.Conn.Call(nil, &buckets, bucketsUri)
	if callErr != nil {
		if v, ok := callErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("code: %d, %s, xreqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			err = callErr
		}
	}
	return
}

type BucketDomain struct {
	Domain string `json:"domain"`
	Owner  int    `json:"owner"`
}

func GetDomainsOfBucket(mac *digest.Mac, bucket string) (domains []BucketDomain, err error) {
	domains = make([]BucketDomain, 0)
	client := rs.NewMac(mac)
	getDomainsUrl := fmt.Sprintf("%s/v7/domain/list?tbl=%s", BUCKET_API_HOST, bucket)
	callErr := client.Conn.Call(nil, &domains, getDomainsUrl)
	if callErr != nil {
		if v, ok := callErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("code: %d, %s, xreqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			err = callErr
		}
	}
	return
}
