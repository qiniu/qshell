package qshell

import (
	"errors"
	"qiniu/api.v6/rs"
)

func BatchRefresh(client *rs.Client, urls []string) (err error) {
	if len(urls) == 0 || len(urls) > 100 {
		err = errors.New("url count invalid, should between [1, 100]")
		return
	}

	postUrl := "http://fusion.qiniuapi.com/refresh"

	postData := map[string][]string{
		"urls": urls,
	}

	err = client.Conn.CallWithForm(nil, nil, postUrl, postData)
	return
}
