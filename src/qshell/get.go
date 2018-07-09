package qshell

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	"strings"
)

// GetFileFromBucket
// check domains
// check src domains
// proxy download
func GetFileFromBucket(mac *digest.Mac, bucket, key, localFile string) (err error) {
	var localFp *os.File

	if localFile == "stdout" {
		localFp = os.Stdout
	} else {
		localFp, err = os.Create(localFile)
		if err != nil {
			return
		}
	}

	defer localFp.Close()

	bWriter := bufio.NewWriter(localFp)
	defer bWriter.Flush()

	var domainOfBucket string

	//get domains of bucket
	domainsOfBucket, gErr := GetDomainsOfBucket(mac, bucket)
	if gErr != nil {
		err = gErr
		return
	}

	if len(domainsOfBucket) == 0 {
		err = fmt.Errorf("Err: no domains found for bucket %s", bucket)
		return
	}

	for _, d := range domainsOfBucket {
		if !strings.HasPrefix(d.Domain, ".") {
			domainOfBucket = d.Domain
			break
		}
	}

	//set proxy
	ioProxyAddress := conf.IO_HOST
	ioProxyAddress = strings.TrimPrefix(ioProxyAddress, "http://")
	ioProxyAddress = strings.TrimPrefix(ioProxyAddress, "https://")

	finalUrl := makePrivateDownloadLink(mac, domainOfBucket, ioProxyAddress, key)
	req, newErr := http.NewRequest("GET", finalUrl, nil)
	if newErr != nil {
		err = newErr
		return
	}
	req.Host = domainOfBucket
	resp, getErr := http.DefaultClient.Do(req)
	if getErr != nil {
		err = getErr
		return
	}
	defer resp.Body.Close()
	_, cpErr := io.Copy(bWriter, resp.Body)
	if cpErr != nil {
		err = cpErr
		return
	}

	return
}
