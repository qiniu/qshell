package download

import (
	"fmt"
	"net/http"

	"github.com/qiniu/qshell/v2/iqshell/common/client"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

type DownloadApiInfo struct {
	downloadUrl    string
	Bucket         string
	Key            string
	IsPublicBucket bool
	UseGetFileApi  bool
	Host           string
	Referer        string
	RangeFromBytes int64
	RangeToBytes   int64
	CheckSize      bool
	FileSize       int64
	CheckHash      bool
	FileHash       string
	Progress       progress.Progress
}

type downloaderFile struct{}

func (d *downloaderFile) Download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError) {
	for times := 0; times < 2; times++ {
		if url, e := createDownloadUrl(info); e != nil {
			return nil, e
		} else {
			info.downloadUrl = url
		}
		response, err = d.download(info)
		log.DebugF("Simple Download[%d] %s, err:%+v", times, info.downloadUrl, err)
		if err == nil {
			break
		}

		if response == nil {
			continue
		}

		if (response.StatusCode > 399 && response.StatusCode < 500) ||
			response.StatusCode == 612 || response.StatusCode == 631 {
			log.DebugF("Simple Stop download %s, because %+v", info.downloadUrl, err)
			break
		}
	}
	return
}

func (d *downloaderFile) download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError) {
	headers := http.Header{}
	if len(info.Host) > 0 {
		headers.Add("Host", info.Host)
	}

	// 设置断点续传
	if info.RangeFromBytes >= 0 && info.RangeToBytes >= 0 && (info.RangeFromBytes+info.RangeToBytes) > 0 {
		if info.RangeFromBytes > 0 && info.RangeFromBytes == info.RangeToBytes {
			return &http.Response{
				Status:     "already download",
				StatusCode: 200,
			}, nil
		} else if info.RangeToBytes == 0 {
			headers.Add("Range", fmt.Sprintf("bytes=%d-", info.RangeFromBytes))
		} else {
			headers.Add("Range", fmt.Sprintf("bytes=%d-%d", info.RangeFromBytes, info.RangeToBytes))
		}
	}

	// 配置 referer
	if len(info.Referer) > 0 {
		headers.Add("Referer", info.Referer)
	}

	if workspace.IsCmdInterrupt() {
		return nil, data.CancelError
	}
	response, rErr := client.DefaultStorageClient().DoRequest(workspace.GetContext(), "GET", info.downloadUrl, headers)
	if info.CheckHash && len(info.FileHash) != 0 && response != nil && response.Header != nil {
		etag := response.Header.Get("Etag")
		etag = utils.ParseEtag(etag)
		if len(etag) > 0 && etag != info.FileHash {
			return nil, data.NewEmptyError().AppendDescF("file has change, hash before:%s now:%s", info.FileHash, etag)
		}
	}

	if rErr == nil {
		return response, nil
	}

	cErr := data.ConvertError(rErr)
	if cErr.Code <= 0 {
		cErr.Code = data.ErrorCodeUnknown
	}
	return response, cErr
}
