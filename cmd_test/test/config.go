package test

import "os"

const (
	Bucket                 = "qshell-na0"
	BucketDomain           = "qshell-na0.qiniupkg.com"
	BucketObjectListString = `
https://qshell-na0.qiniupkg.com/hello1.json
https://qshell-na0.qiniupkg.com/hello2.json
https://qshell-na0.qiniupkg.com/hello3.json
https://qshell-na0.qiniupkg.com/hello4.json
https://qshell-na0.qiniupkg.com/hello5.json
https://qshell-na0.qiniupkg.com/hello6.json
https://qshell-na0.qiniupkg.com/hello7.json
`
)

var (
	Debug = true

	AccessKey = os.Getenv("accessKey")
	SecretKey = os.Getenv("secretKey")
)
