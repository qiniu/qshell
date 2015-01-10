package qshell

import (
	"fmt"
	"github.com/qiniu/api/rs"
)

func PrintStat(bucket string, key string, entry rs.Entry) {
	statInfo := fmt.Sprintf("%-20s%-20s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "Fsize:", entry.Fsize)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "PutTime:", entry.PutTime)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "MimeType:", entry.MimeType)
	fmt.Println(statInfo)
}
