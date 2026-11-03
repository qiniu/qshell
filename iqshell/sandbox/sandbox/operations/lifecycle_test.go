package operations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func TestApplyLifecycleParams(t *testing.T) {
	params := sandbox.CreateParams{}
	require.NoError(t, applyLifecycleParams(CreateInfo{}, &params))
	require.Nil(t, params.AutoResume)
	require.Nil(t, params.AutoPauseMemory)
	disabled, enabled := false, true
	require.NoError(t, applyLifecycleParams(CreateInfo{AutoPauseMemory: &disabled, AutoResume: &disabled}, &params))
	require.NotNil(t, params.AutoPauseMemory)
	require.False(t, *params.AutoPauseMemory)
	require.False(t, params.AutoResume.Enabled)
	require.Error(t, applyLifecycleParams(CreateInfo{AutoPauseMemory: &disabled, AutoResume: &enabled}, &params))
}

func TestPauseSendsFilesystemOnlyWithoutConnecting(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/sandboxes/sb-1/pause", r.URL.Path)
		var body map[string]any
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
			assert.Equal(t, false, body["memory"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	t.Setenv("QINIU_API_KEY", "test-key")
	t.Setenv("QINIU_SANDBOX_API_URL", server.URL)
	disabled := false
	Pause(PauseInfo{SandboxIDs: []string{"sb-1"}, Memory: &disabled})
	require.True(t, called)
}

func TestInvalidLifecycleSetsFailureStatus(t *testing.T) {
	data.SetCmdStatus(data.StatusOK)
	t.Cleanup(func() { data.SetCmdStatus(data.StatusOK) })
	t.Setenv("QINIU_API_KEY", "test-key")
	disabled, enabled := false, true
	Create(CreateInfo{TemplateID: "base", Detach: true, AutoPauseMemory: &disabled, AutoResume: &enabled})
	require.NotEqual(t, data.StatusOK, data.GetCmdStatus())
}
