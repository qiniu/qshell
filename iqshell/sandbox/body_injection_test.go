package sandbox

import (
	"testing"

	sdk "github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/require"
)

func TestHTTPBodyFieldsPreserveJSONStrings(t *testing.T) {
	parts := InjectionParts{HTTP: &sdk.HTTPInjection{BaseURL: "https://example.com"}}
	require.NoError(t, ApplyHTTPBodyFields(&parts, `{"api_key":"a,b=c;d","empty":""}`, `{"model":"demo"}`))
	require.Equal(t, "a,b=c;d", (*parts.HTTP.Body)["api_key"])
	require.Contains(t, *parts.HTTP.Body, "empty")
	require.Equal(t, "demo", (*parts.HTTP.IfBody)["model"])
	for _, invalid := range []string{`null`, `[]`, `{"secret":123}`, `{"secret":null}`, `{"secret":`} {
		err := ApplyHTTPBodyFields(&parts, invalid, "")
		require.Error(t, err)
		require.NotContains(t, err.Error(), invalid)
	}
	require.Error(t, ApplyHTTPBodyFields(&InjectionParts{}, `{"key":"v"}`, ""))
}
