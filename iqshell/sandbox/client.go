package sandbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/subosito/gotenv"

	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
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
	response, err := t.base.RoundTrip(req)
	if response != nil && response.StatusCode == http.StatusOK {
		if token, ok := req.Context().Value(listPageTokenKey{}).(*listPageToken); ok {
			token.value = response.Header.Get("X-Next-Token")
		}
	}
	return response, err
}

// Environment variable names for sandbox configuration.
const (
	// Qiniu-specific environment variables (highest priority).
	EnvQiniuSandboxAPIURL = "QINIU_SANDBOX_API_URL"
	EnvQiniuAPIKey        = "QINIU_API_KEY"

	// E2B-compatible environment variables (fallback).
	EnvE2BAPIURL = "E2B_API_URL"
	EnvE2BAPIKey = "E2B_API_KEY"

	// Qiniu AK/SK environment variables (fallback when no qshell account is configured).
	EnvQiniuAccessKey = "QINIU_ACCESS_KEY"
	EnvQiniuSecretKey = "QINIU_SECRET_KEY"
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

// NewSandboxClient creates a sandbox client for sandbox runtime / template operations.
// API Key is required; AK/SK Credentials are included when available.
//
// Credential sources (each independent):
//   - API Key: env QINIU_API_KEY > E2B_API_KEY (.env file is also consulted).
//   - AK/SK:   qshell active account (~/.qshell) > env QINIU_ACCESS_KEY/QINIU_SECRET_KEY.
//
// Returns an error if API Key is missing — sandbox runtime / template endpoints
// only accept API Key authentication.
func NewSandboxClient() (*sandbox.Client, error) {
	return NewSandboxClientWithRetryMax(nil)
}

// NewSandboxClientWithRetryMax 使用可选的重试上限创建沙箱客户端。
// retryMax 为 nil 时保留 SDK 默认行为，包括读取 SANDBOX_RETRY_MAX。
func NewSandboxClientWithRetryMax(retryMax *int) (*sandbox.Client, error) {
	loadDotEnv()

	apiKey, endpoint := resolveConfig()
	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured for sandbox/template operations: please set %s or %s environment variable", EnvQiniuAPIKey, EnvE2BAPIKey)
	}

	return buildSandboxClient(apiKey, endpoint, resolveCredentials(), retryMax)
}

// NewInjectionRuleClient creates a sandbox client for injection rule operations.
// AK/SK Credentials are required; API Key is included when available.
//
// Credential sources (see NewSandboxClient for details).
//
// Returns an error if AK/SK is missing — injection rule endpoints sign requests
// with Qiniu credentials and cannot fall back to API Key.
func NewInjectionRuleClient() (*sandbox.Client, error) {
	loadDotEnv()

	creds := resolveCredentials()
	if creds == nil {
		return nil, fmt.Errorf("qiniu credentials (AK/SK) not configured for injection rule operations: please configure a qshell account via 'qshell user' or set %s/%s environment variables", EnvQiniuAccessKey, EnvQiniuSecretKey)
	}

	apiKey, endpoint := resolveConfig()
	return buildSandboxClient(apiKey, endpoint, creds, nil)
}

// buildSandboxClient constructs the underlying SDK client. apiKey may be empty
// (only meaningful for InjectionRule paths); creds may be nil.
func buildSandboxClient(apiKey, endpoint string, creds *auth.Credentials, retryMax *int) (*sandbox.Client, error) {
	return sandbox.NewClient(&sandbox.Config{
		APIKey:      apiKey,
		Credentials: creds,
		Endpoint:    endpoint,
		RetryMax:    retryMax,
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

// workspaceAccountLookup is the seam used by resolveCredentials to fetch the
// active qshell account. Production points it at workspace.GetAccount; tests
// override it to inject deterministic credentials. Returns empty strings when
// no active account is configured.
var workspaceAccountLookup = func() (accessKey, secretKey string) {
	acc, err := workspace.GetAccount()
	if err != nil {
		return "", ""
	}
	return acc.AccessKey, acc.SecretKey
}

// resolveCredentials returns AK/SK credentials, preferring the active qshell
// account over environment variables. Returns nil when neither source provides
// a complete pair.
//
// Source priority:
//  1. workspace.GetAccount() — the active qshell account loaded from ~/.qshell.
//     Requires the caller (cmd layer) to have invoked iqshell.CheckAndLoad first.
//     Note: we read the account directly instead of workspace.GetConfig().Credentials
//     because Load() rebuilds cfg via Merge with defaultConfig(), which overwrites
//     the credentials populated by loadUserInfo().
//  2. QINIU_ACCESS_KEY / QINIU_SECRET_KEY environment variables.
func resolveCredentials() *auth.Credentials {
	ak, sk := workspaceAccountLookup()
	if ak == "" || sk == "" {
		ak = os.Getenv(EnvQiniuAccessKey)
		sk = os.Getenv(EnvQiniuSecretKey)
	}
	if ak == "" || sk == "" {
		return nil
	}
	return &auth.Credentials{AccessKey: ak, SecretKey: []byte(sk)}
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
