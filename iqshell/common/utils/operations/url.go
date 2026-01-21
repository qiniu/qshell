package operations

import (
	"net/url"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type UrlInfo struct {
	Url string
}

func (info *UrlInfo) Check() *data.CodeError {
	if len(info.Url) == 0 {
		return alert.CannotEmptyError("Data", "")
	}
	return nil
}

func UrlEncode(cfg *iqshell.Config, info UrlInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	dataEncoded := url.PathEscape(info.Url)
	log.Alert(dataEncoded)
}

func UrlDecode(cfg *iqshell.Config, info UrlInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	dataDecoded, err := url.PathUnescape(info.Url)
	if err != nil {
		data.SetCmdStatusError()
		log.Error("Failed to unescape data `", info.Url, "'")
	} else {
		log.Alert(dataDecoded)
	}
}
