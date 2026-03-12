package sandbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/subosito/gotenv"
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

// loadDotEnv loads variables from the .env file in the current directory.
// Only variables not already set in the OS environment are loaded (OS takes priority).
// Missing or unreadable .env files are silently ignored.
func loadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	env, err := gotenv.StrictParse(f)
	if err != nil {
		return
	}
	for key, value := range env {
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, strings.TrimSpace(value))
		}
	}
}

// NewSandboxClient creates a new sandbox client by reading configuration from environment variables.
// It first loads .env file from the current directory (OS env vars take priority).
// Priority: QINIU_SANDBOX_API_URL / QINIU_API_KEY > E2B_API_URL / E2B_API_KEY.
func NewSandboxClient() (*sandbox.Client, error) {
	loadDotEnv()

	apiKey, endpoint := resolveConfig()
	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured, please set %s or %s environment variable", EnvQiniuAPIKey, EnvE2BAPIKey)
	}

	return sandbox.NewClient(&sandbox.Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
		HTTPClient: &http.Client{
			Transport: &keepaliveTransport{base: http.DefaultTransport},
		},
	})
}

// resolveConfig returns the resolved API key and endpoint from environment variables.
func resolveConfig() (apiKey, endpoint string) {
	apiKey = os.Getenv(EnvQiniuAPIKey)
	if apiKey == "" {
		apiKey = os.Getenv(EnvE2BAPIKey)
	}
	endpoint = os.Getenv(EnvQiniuSandboxAPIURL)
	if endpoint == "" {
		endpoint = os.Getenv(EnvE2BAPIURL)
	}
	if endpoint == "" {
		endpoint = sandbox.DefaultEndpoint
	}
	return apiKey, endpoint
}

// ResumeSandbox resumes a paused sandbox by calling POST /sandboxes/{id}/resume.
// The SDK Client does not expose a Resume method, so we call the API directly.
func ResumeSandbox(sandboxID string, timeout *int32) error {
	loadDotEnv()

	apiKey, endpoint := resolveConfig()
	if apiKey == "" {
		return fmt.Errorf("API key not configured, please set %s or %s environment variable", EnvQiniuAPIKey, EnvE2BAPIKey)
	}

	body := map[string]any{}
	if timeout != nil {
		body["timeout"] = *timeout
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/sandboxes/%s/resume", strings.TrimRight(endpoint, "/"), sandboxID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resume request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("api error: status %d, body: %s", resp.StatusCode, string(respBody))
}
