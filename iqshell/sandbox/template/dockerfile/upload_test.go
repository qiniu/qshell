package dockerfile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeFilesHash(t *testing.T) {
	// 创建包含测试文件的临时目录
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0644))

	hash1, err := ComputeFilesHash(".", "/app/", dir, nil)
	require.NoError(t, err)
	assert.Len(t, hash1, 64) // SHA-256 十六进制

	// 相同文件产生相同哈希
	hash2, err := ComputeFilesHash(".", "/app/", dir, nil)
	require.NoError(t, err)
	assert.Equal(t, hash1, hash2)

	// 不同目标路径产生不同哈希
	hash3, err := ComputeFilesHash(".", "/other/", dir, nil)
	require.NoError(t, err)
	assert.NotEqual(t, hash1, hash3)
}

func TestComputeFilesHash_SingleFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.py"), []byte("print('hello')"), 0644))

	hash, err := ComputeFilesHash("app.py", "/app/", dir, nil)
	require.NoError(t, err)
	assert.Len(t, hash, 64)
}

func TestComputeFilesHash_WithIgnore(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.py"), []byte("code"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "test.log"), []byte("log data"), 0644))

	hashWithLog, err := ComputeFilesHash(".", "/app/", dir, nil)
	require.NoError(t, err)

	hashNoLog, err := ComputeFilesHash(".", "/app/", dir, []string{"*.log"})
	require.NoError(t, err)

	assert.NotEqual(t, hashWithLog, hashNoLog)
}

func TestCollectAndUpload(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644))

	var received bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/gzip", r.Header.Get("Content-Type"))
		received = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := CollectAndUpload(context.Background(), server.URL, ".", dir, nil)
	require.NoError(t, err)
	assert.True(t, received)
}

func TestReadDockerignore(t *testing.T) {
	dir := t.TempDir()

	// 无 .dockerignore 文件
	patterns := ReadDockerignore(dir)
	assert.Nil(t, patterns)

	// 创建 .dockerignore
	content := "*.log\n# comment\n\nnode_modules\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".dockerignore"), []byte(content), 0644))

	patterns = ReadDockerignore(dir)
	assert.Equal(t, []string{"*.log", "node_modules"}, patterns)
}

func TestIsIgnored(t *testing.T) {
	tests := []struct {
		name     string
		relPath  string
		patterns []string
		want     bool
	}{
		{"无规则", "file.txt", nil, false},
		{"匹配规则", "test.log", []string{"*.log"}, true},
		{"不匹配规则", "app.py", []string{"*.log"}, false},
		{"取反规则", "important.log", []string{"*.log", "!important.log"}, true},
		{"基名匹配", "dir/file.log", []string{"*.log"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIgnored(tt.relPath, tt.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}
