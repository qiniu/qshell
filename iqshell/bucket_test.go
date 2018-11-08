package qshell

import (
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"os"
	"testing"
)

func TestGetBucketInfo(t *testing.T) {
	ak := os.Getenv("AccessKey")
	sk := os.Getenv("SecretKey")
	bucket := os.Getenv("Bucket")
	mac := digest.Mac{ak, []byte(sk)}
	bucketInfo, gErr := GetBucketInfo(&mac, bucket)
	if gErr != nil {
		t.Fatal(gErr)
	}
	t.Log(bucketInfo.Region)
}
