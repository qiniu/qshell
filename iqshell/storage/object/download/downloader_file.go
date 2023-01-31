package download

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/http"
)

type DownloadApiInfo struct {
	Url            string
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

type downloaderFile struct {
}

func (d *downloaderFile) Download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError) {
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
	response, rErr := defaultClient.DoRequest(workspace.GetContext(), "GET", info.Url, headers)
	if info.CheckHash && len(info.FileHash) != 0 && response != nil && response.Header != nil {
		etag := fmt.Sprintf(response.Header.Get("Etag"))
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
