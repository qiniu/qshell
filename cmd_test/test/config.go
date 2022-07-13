package test

import (
	"fmt"
	"os"
	"strings"
)

var (
	Debug          = true
	Local          = len(os.Getenv("QSHELL_LOCAL")) > 0
	ShouldTestUser = len(os.Getenv("QSHELL_Test_User")) > 0

	AccessKey      = os.Getenv("accessKey")
	SecretKey      = os.Getenv("secretKey")
	Bucket         = testBucket()
	BucketDomain   = testBucketDomain()
	UploadDomain   = testUploadDomain()
	DocumentOption = "--doc"
)

func testBucket() string {
	if b := os.Getenv("bucket"); len(b) > 0 {
		return b
	} else {
		return "qshell-z1"
	}
}

func testBucketDomain() string {
	if b := os.Getenv("bucketDomain"); len(b) > 0 {
		return b
	} else {
		return "qshell-z1.qiniupkg.com"
	}
}

func testUploadDomain() string {
	if b := os.Getenv("uploadDomain"); len(b) > 0 {
		return b
	} else {
		return "up-z1.qiniup.com"
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
