package qshell

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rsf"
	"github.com/qiniu/log"
	"io"
	"os"
)

type ListBucket struct {
	Account
}

func (this *ListBucket) List(bucket string, prefix string, listResultFile string) {
	fp, openErr := os.OpenFile(listResultFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if openErr != nil {
		log.Error(fmt.Sprintf("Failed to open list result file `%s'", listResultFile))
		return
	}
	defer fp.Close()
	bw := bufio.NewWriter(fp)

	mac := digest.Mac{this.AccessKey, []byte(this.SecretKey)}
	client := rsf.New(&mac)
	marker := ""
	limit := 1000
	run := true
	for run {
		entries, markerOut, err := client.ListPrefix(nil, bucket, prefix, marker, limit)
		if err != nil {
			if err == io.EOF {
				run = false
			}
		} else {
			if markerOut == "" {
				run = false
			} else {
				marker = markerOut
			}
		}
		//append entries
		for _, entry := range entries {
			lineData := fmt.Sprintf("%s\t%d\r\n", entry.Key, entry.Fsize)
			_, wErr := bw.WriteString(lineData)
			if wErr != nil {
				log.Error(fmt.Sprintf("Write line data `%s' to list result file failed.", lineData))
			}
			fErr := bw.Flush()
			if fErr != nil {
				log.Error("Flush data to list result file error", err)
			}
		}
	}
}
