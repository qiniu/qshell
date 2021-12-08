package tools

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"net/url"
)

type UrlInfo struct {
	Url string
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