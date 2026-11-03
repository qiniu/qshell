package operations

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func TestBodyInjectionAndFailureStatus(t *testing.T) {
	spec, err := buildInjectionSpec(injectionInput{Type: "http", BaseURL: "https://example.com", Body: `{"api_key":"fake,secret"}`, IfBody: `{"model":"demo"}`})
	require.NoError(t, err)
	require.Equal(t, "fake,secret", (*spec.HTTP.Body)["api_key"])
	require.True(t, shouldUpdateInjection(injectionInput{Body: `{}`}))
	require.Contains(t, formatInjectionConditions(spec), "model")
	require.NotContains(t, formatInjectionConditions(spec), "demo")
	data.SetCmdStatus(data.StatusOK)
	t.Cleanup(func() { data.SetCmdStatus(data.StatusOK) })
	Create(CreateInfo{Name: "invalid", Type: "http", BaseURL: "https://example.com", Body: `{"api_key":null}`})
	require.NotEqual(t, data.StatusOK, data.GetCmdStatus())
}
