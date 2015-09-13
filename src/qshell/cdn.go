package qshell

import (
	"errors"
	"github.com/qiniu/api.v6/rs"
)

func BatchRefresh(client *rs.Client, urls []string) (err error) {
	if len(urls) == 0 || len(urls) > 10 {
		err = errors.New("url count invalid, should between [1, 10]")
		return
	}

	postUrl := "http://cdnmgr.qbox.me:15001/refresh/"

	postData := map[string][]string{
		"urls": urls,
	}

	err = client.Conn.CallWithForm(nil, nil, postUrl, postData)
	return
}
