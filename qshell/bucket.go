package qshell

import (
	"fmt"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"github.com/tonycai653/iqshell/qiniu/rpc"
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
	bucketUri := fmt.Sprintf("%s/bucket/%s", BUCKET_RS_HOST, bucket)
	callErr := client.Conn.Call(nil, &bucketInfo, bucketUri)
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

func GetDomainsOfBucket(mac *digest.Mac, bucket string) (domains []string, err error) {
	domains = make([]string, 0)
	client := rs.NewMac(mac)
	getDomainsUrl := fmt.Sprintf("%s/v6/domain/list", BUCKET_API_HOST)
	postData := map[string][]string{
		"tbl": []string{bucket},
	}
	callErr := client.Conn.CallWithForm(nil, &domains, getDomainsUrl, postData)
	if callErr != nil {
		if v, ok := callErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("code: %d, %s, xreqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			err = callErr
		}
	}
	return
}
