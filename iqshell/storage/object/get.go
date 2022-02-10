package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GetApiInfo struct {
	Bucket     string
	Key        string
	Domain     string
	SaveToFile string
}

func GetObject(info GetApiInfo) (err error) {

	if len(info.Bucket) == 0 {
		return errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.Key) == 0 {
		return errors.New(alert.CannotEmpty("key", ""))
	}

	if len(info.SaveToFile) == 0 {
		info.SaveToFile = info.Key
	}

	data, err := getObjectFromServer(info)
	if err != nil {
		return err
	}

	err = writeDataToFile(data, info.SaveToFile)
	data.Close()
	return err
}

func getObjectFromServer(info GetApiInfo) (io.ReadCloser, error) {
	entryUri := strings.Join([]string{info.Bucket, info.Key}, ":")

	var reqHost string
	reqHost = workspace.GetConfig().Hosts.GetOneRs()
	if reqHost == "" {
		bucketManager, err := bucket.GetBucketManager()
		if err != nil {
			return nil, err
		}

		zone, err := bucketManager.Zone(info.Bucket)
		if err != nil {
			return nil, err
		}
		reqHost = zone.RsHost
	}

	if !strings.HasPrefix(reqHost, "http") {
		reqHost = "http://" + reqHost
	}
	url := strings.Join([]string{reqHost, "get", utils.Encode(entryUri)}, "/")

	mac, err := workspace.GetMac()
	if err != nil {
		return nil, err
	}

	var data struct {
		URL string `json:"url"`
	}
	ctx := auth.WithCredentials(workspace.GetContext(), mac)
	headers := http.Header{}

	err = storage.DefaultClient.Call(ctx, &data, "GET", url, headers)
	if err != nil {
		return nil, err
	}

	resp, err := storage.DefaultClient.DoRequest(workspace.GetContext(), "GET", data.URL, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("Qget: http respcode: %d, respbody: %s\n", resp.StatusCode, string(body)))
	}

	return resp.Body, nil
}

func writeDataToFile(data io.Reader, destFile string) (err error) {
	var file *os.File
	for {
		file, err = os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err == nil {
			break
		} else if os.IsNotExist(err) {
			destDir := filepath.Dir(destFile)
			if err = os.MkdirAll(destDir, 0700); err != nil {
				err = errors.New("Qget: err:" + err.Error())
				break
			}
		}
	}

	if err != nil {
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, data); err != nil {
		err = errors.New("Qget: err:" + err.Error())
	}

	return
}
