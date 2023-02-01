package test

import (
	"fmt"
	"os"
	"strings"
)

var (
	AccessKey      = os.Getenv("accessKey")
	SecretKey      = os.Getenv("secretKey")
	Bucket         = testBucket()
	BucketDomain   = testBucketDomain()
	UploadDomain   = testUploadDomain()
	DocumentOption = "--doc"
)

func IsDebug() bool {
	return true
}

func testBucket() string {
	if b := os.Getenv("bucket"); len(b) > 0 {
		return b
	} else {
		return "qshell-z0"
	}
}

func testBucketDomain() string {
	if b := os.Getenv("bucketDomain"); len(b) > 0 {
		return b
	} else {
		return "qshell-z0-src.qiniupkg.com"
	}
}

func testUploadDomain() string {
	if b := os.Getenv("uploadDomain"); len(b) > 0 {
		return b
	} else {
		return "up-z0.qiniup.com"
	}
}

var (
	BucketNotExist      = "qshell-na0-mock"
	BucketObjectDomain  = fmt.Sprintf("https://%s/hello1_test.json", BucketDomain)
	BucketObjectDomains = []string{
		fmt.Sprintf("https://%s/hello1_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello2_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello3_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello4_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello5_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello6_test.json", BucketDomain),
		fmt.Sprintf("https://%s/hello7_test.json", BucketDomain),
	}
	BucketObjectDomainsString = strings.ReplaceAll(`https://domain/hello1_test.json
https://domain/hello2_test.json
https://domain/hello3_test.json
https://domain/hello4_test.json
https://domain/hello5_test.json
https://domain/hello6_test.json
https://domain/hello7_test.json
`, "domain", BucketDomain)

	Key         = "hello1_test.json"
	ImageKey    = "image.png"
	KeyNotExist = "hello_mock.json"
	OriginKeys  = []string{"hello1.json", "hello2.json", "hello3.json", "hello4.json", "hello5.json", "hello6.json", "hello7.json"}
	Keys        = []string{"hello1_test.json", "hello2_test.json", "hello3_test.json", "hello4_test.json", "hello5_test.json", "hello6_test.json", "hello7_test.json"}
	KeysString  = `hello1_test.json
hello2_test.json
hello3_test.json
hello4_test.json
hello5_test.json
hello6_test.json
hello7_test.json`
)
