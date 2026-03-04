package sandbox

import (
	"fmt"
	"os"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

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
	})
}
