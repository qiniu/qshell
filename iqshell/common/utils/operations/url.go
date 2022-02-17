package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"net/url"
)

type UrlInfo struct {
	Url string
}

func (info *UrlInfo) Check() error {
	if len(info.Url) == 0 {
		return alert.CannotEmptyError("Data", "")
	}
	return nil
}

func UrlEncode(info UrlInfo) {
	dataEncoded := url.PathEscape(info.Url)
	log.Alert(dataEncoded)
}

func UrlDecode(info UrlInfo) {
	dataDecoded, err := url.PathUnescape(info.Url)
	if err != nil {
		log.Error("Failed to unescape data `", info.Url, "'")
	} else {
		log.Alert(dataDecoded)
	}
}
