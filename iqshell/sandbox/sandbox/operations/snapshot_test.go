package operations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotOperationsRequests(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			assert.Equal(t, "/sandboxes/sb-1/snapshots", r.URL.Path)
			assert.Equal(t, "capture-1", r.Header.Get("Idempotency-Key"))
			var body map[string]any
			if !assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
				w.WriteHeader(400)
				return
			}
			assert.Equal(t, false, body["memory"])
			assert.Equal(t, "demo:v1", body["name"])
			w.WriteHeader(201)
			_, _ = w.Write([]byte(`{"snapshotID":"snap-1","names":["demo:v1"]}`))
			return
		}
		assert.Equal(t, "/snapshots", r.URL.Path)
		assert.Equal(t, "cursor", r.URL.Query().Get("nextToken"))
		assert.Equal(t, "demo", r.URL.Query().Get("name"))
		assert.Equal(t, "sb-1", r.URL.Query().Get("sandboxID"))
		assert.Equal(t, "2", r.URL.Query().Get("limit"))
		w.Header().Set("X-Next-Token", "next")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()
	t.Setenv("QINIU_API_KEY", "test-key")
	t.Setenv("QINIU_SANDBOX_API_URL", server.URL)
	require.NoError(t, SnapshotCreate(t.Context(), SnapshotCreateInfo{SandboxID: "sb-1", Name: "demo:v1", Memory: false, IdempotencyKey: "capture-1", Format: "json"}))
	require.NoError(t, SnapshotList(t.Context(), SnapshotListInfo{SandboxID: "sb-1", Name: "demo", NextToken: "cursor", Limit: 2, Format: "json"}))
	require.Equal(t, 2, calls)
	require.Error(t, SnapshotCreate(t.Context(), SnapshotCreateInfo{Format: "json"}))
	require.Error(t, SnapshotList(t.Context(), SnapshotListInfo{Limit: -1, Format: "json"}))
	require.Equal(t, 2, calls, "无效参数不应发送请求")
}

func TestSnapshotOperationPropagatesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"code":403,"message":"not allowed"}`))
	}))
	defer server.Close()
	t.Setenv("QINIU_API_KEY", "test-key")
	t.Setenv("QINIU_SANDBOX_API_URL", server.URL)
	err := SnapshotCreate(t.Context(), SnapshotCreateInfo{SandboxID: "sb-1", IdempotencyKey: "retry-key", Memory: true, Format: "json"})
	require.ErrorContains(t, err, "retry-key")
	require.ErrorContains(t, err, "403")
}
