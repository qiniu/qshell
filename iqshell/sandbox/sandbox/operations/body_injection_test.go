package operations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInlineHTTPBodyJSON(t *testing.T) {
	injection, err := parseInlineSandboxInjection(`type=http,base-url=https://example.com,body={"api_key":"a,b=c;d","another":"v"},if-body={"model":"demo"}`)
	require.NoError(t, err)
	require.NotNil(t, injection.HTTP)
	require.NotNil(t, injection.HTTP.Body)
	require.Equal(t, "a,b=c;d", (*injection.HTTP.Body)["api_key"])
	require.Equal(t, "v", (*injection.HTTP.Body)["another"])
	require.Equal(t, "demo", (*injection.HTTP.IfBody)["model"])
}
