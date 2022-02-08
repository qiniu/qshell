package test

import "os"

var (
	Bucket              = "qshell-na0"
	BucketDomain        = "qshell-na0.qiniupkg.com"
	BucketObjectDomain  = "https://qshell-na0.qiniupkg.com/hello1.json"
	BucketObjectDomains = []string{
		"https://qshell-na0.qiniupkg.com/hello1.json",
		"https://qshell-na0.qiniupkg.com/hello2.json",
		"https://qshell-na0.qiniupkg.com/hello3.json",
		"https://qshell-na0.qiniupkg.com/hello4.json",
		"https://qshell-na0.qiniupkg.com/hello5.json",
		"https://qshell-na0.qiniupkg.com/hello6.json",
		"https://qshell-na0.qiniupkg.com/hello7.json",
	}
	BucketObjectDomainsString = `
https://qshell-na0.qiniupkg.com/hello1.json
https://qshell-na0.qiniupkg.com/hello2.json
https://qshell-na0.qiniupkg.com/hello3.json
https://qshell-na0.qiniupkg.com/hello4.json
https://qshell-na0.qiniupkg.com/hello5.json
https://qshell-na0.qiniupkg.com/hello6.json
https://qshell-na0.qiniupkg.com/hello7.json
`
	Key        = "hello1.json"
	Keys       = []string{"hello1.json", "hello2.json", "hello3.json", "hello4.json", "hello5.json", "hello6.json", "hello7.json"}
	KeysString = `
hello1.json
hello2.json
hello3.json
hello4.json
hello5.json
hello6.json
hello7.json
`
)

var (
	Debug = true

	AccessKey = os.Getenv("accessKey")
	SecretKey = os.Getenv("secretKey")
)
