package sandbox

import (
	"os"
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
)

// clearEnvVars unsets all sandbox-related env vars and returns a cleanup function.
func clearEnvVars(t *testing.T) {
	t.Helper()
	envs := []string{EnvQiniuAPIKey, EnvE2BAPIKey, EnvQiniuSandboxAPIURL, EnvE2BAPIURL}
	saved := make(map[string]string)
	for _, k := range envs {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range envs {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

func TestResolveConfig_QiniuPriority(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "qiniu-key")
	os.Setenv(EnvE2BAPIKey, "e2b-key")
	os.Setenv(EnvQiniuSandboxAPIURL, "https://qiniu.example.com")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "qiniu-key", apiKey)
	assert.Equal(t, "https://qiniu.example.com", endpoint)
}

func TestResolveConfig_FallbackToE2B(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvE2BAPIKey, "e2b-key")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "e2b-key", apiKey)
	assert.Equal(t, "https://e2b.example.com", endpoint)
}

func TestResolveConfig_DefaultEndpoint(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "some-key")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "some-key", apiKey)
	assert.Equal(t, sandbox.DefaultEndpoint, endpoint)
}

func TestResolveConfig_AllEmpty(t *testing.T) {
	clearEnvVars(t)

	apiKey, endpoint := resolveConfig()
	assert.Empty(t, apiKey)
	assert.Equal(t, sandbox.DefaultEndpoint, endpoint)
}

func TestResolveConfig_QiniuKeyWithE2BEndpoint(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "qiniu-key")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "qiniu-key", apiKey)
	assert.Equal(t, "https://e2b.example.com", endpoint)
}
