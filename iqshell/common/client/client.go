package client

import (
	"net"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
)

var defaultClient = storage.Client{
	Client: &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          4000,
			MaxIdleConnsPerHost:   1000,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
}

func DefaultStorageClient() storage.Client {
	return defaultClient
}
