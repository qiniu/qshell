package sandbox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

// keepalivePingIntervalSec matches the JS SDK's KEEPALIVE_PING_INTERVAL_SEC (50s).
// This header tells the envd server to send periodic keepalive pings on gRPC streams,
// preventing proxies/load balancers from closing idle connections.
const keepalivePingIntervalSec = "50"

// keepalivePingHeader is the HTTP header name for the keepalive ping interval.
const keepalivePingHeader = "Keepalive-Ping-Interval"

// keepaliveTransport wraps an http.RoundTripper to inject the Keepalive-Ping-Interval header.
type keepaliveTransport struct {
	base http.RoundTripper
}

func (t *keepaliveTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(keepalivePingHeader, keepalivePingIntervalSec)
	return t.base.RoundTrip(req)
}

// Environment variable names for sandbox configuration.
const (
	// Qiniu-specific environment variables (highest priority).
	EnvQiniuSandboxAPIURL = "QINIU_SANDBOX_API_URL"
	EnvQiniuAPIKey        = "QINIU_API_KEY"

	// E2B-compatible environment variables (fallback).
	EnvE2BAPIURL = "E2B_API_URL"
	EnvE2BAPIKey = "E2B_API_KEY"
)

// NewSandboxClient creates a new sandbox client by reading configuration from environment variables.
// Priority: QINIU_SANDBOX_API_URL / QINIU_API_KEY > E2B_API_URL / E2B_API_KEY.
func NewSandboxClient() (*sandbox.Client, error) {
	apiKey := os.Getenv(EnvQiniuAPIKey)
	if apiKey == "" {
		apiKey = os.Getenv(EnvE2BAPIKey)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured, please set %s or %s environment variable", EnvQiniuAPIKey, EnvE2BAPIKey)
	}

	endpoint := os.Getenv(EnvQiniuSandboxAPIURL)
	if endpoint == "" {
		endpoint = os.Getenv(EnvE2BAPIURL)
	}

	return sandbox.NewClient(&sandbox.Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
		HTTPClient: &http.Client{
			Transport: &keepaliveTransport{base: http.DefaultTransport},
		},
	})
}
