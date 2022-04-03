package test

import "os"

var (
	Bucket              = "qshell-na0"
	BucketNotExist      = "qshell-na0-mock"
	BucketDomain        = "qshell-na0.qiniupkg.com"
	BucketObjectDomain  = "https://qshell-na0.qiniupkg.com/hello1_test.json"
	BucketObjectDomains = []string{
		"https://qshell-na0.qiniupkg.com/hello1_test.json",
		"https://qshell-na0.qiniupkg.com/hello2_test.json",
		"https://qshell-na0.qiniupkg.com/hello3_test.json",
		"https://qshell-na0.qiniupkg.com/hello4_test.json",
		"https://qshell-na0.qiniupkg.com/hello5_test.json",
		"https://qshell-na0.qiniupkg.com/hello6_test.json",
		"https://qshell-na0.qiniupkg.com/hello7_test.json",
	}
	BucketObjectDomainsString = `https://qshell-na0.qiniupkg.com/hello1_test.json
https://qshell-na0.qiniupkg.com/hello2_test.json
https://qshell-na0.qiniupkg.com/hello3_test.json
https://qshell-na0.qiniupkg.com/hello4_test.json
https://qshell-na0.qiniupkg.com/hello5_test.json
https://qshell-na0.qiniupkg.com/hello6_test.json
https://qshell-na0.qiniupkg.com/hello7_test.json
`
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

var (
	Debug = true

	AccessKey = os.Getenv("accessKey")
	SecretKey = os.Getenv("secretKey")

	DocumentOption = "--doc"
)
