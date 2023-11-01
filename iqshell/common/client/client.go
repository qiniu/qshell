package client

import (
	"net"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
)

var defaultClient = storage.Client{
	Client: &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          2000,
			MaxIdleConnsPerHost:   1000,
			ResponseHeaderTimeout: 60 * time.Second,
			IdleConnTimeout:       15 * time.Second,
			TLSHandshakeTimeout:   15 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
}

func DefaultStorageClient() storage.Client {
	return defaultClient
}
