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
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0o644))

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
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.py"), []byte("print('hello')"), 0o644))

	hash, err := ComputeFilesHash("app.py", "/app/", dir, nil)
	require.NoError(t, err)
	assert.Len(t, hash, 64)
}

func TestComputeFilesHash_WithIgnore(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.py"), []byte("code"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "test.log"), []byte("log data"), 0o644))

	hashWithLog, err := ComputeFilesHash(".", "/app/", dir, nil)
	require.NoError(t, err)

	hashNoLog, err := ComputeFilesHash(".", "/app/", dir, []string{"*.log"})
	require.NoError(t, err)

	assert.NotEqual(t, hashWithLog, hashNoLog)
}

func TestCollectAndUpload(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0o644))

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
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".dockerignore"), []byte(content), 0o644))

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
		{"取反规则重新包含", "important.log", []string{"*.log", "!important.log"}, false},
		{"取反规则不影响其他", "other.log", []string{"*.log", "!important.log"}, true},
		{"基名匹配", "dir/file.log", []string{"*.log"}, true},
		// 否定规则: * 排除所有, !README.md 重新包含
		{"全排除后否定包含", "README.md", []string{"*", "!README.md"}, false},
		{"全排除后否定不影响其他", "app.go", []string{"*", "!README.md"}, true},
		// 多重否定: 排除 → 包含 → 再排除
		{"多重否定最终排除", "debug.log", []string{"*.log", "!*.log", "debug.log"}, true},
		{"多重否定最终包含", "info.log", []string{"*.log", "!*.log", "debug.log"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIgnored(tt.relPath, tt.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsWithinContext(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		contextAbs string
		want       bool
	}{
		{"子路径", "/home/user/project/src/main.go", "/home/user/project", true},
		{"相同路径", "/home/user/project", "/home/user/project", true},
		{"逃逸路径", "/home/user/other/file.txt", "/home/user/project", false},
		{"前缀但非子目录", "/home/user/project2/file.txt", "/home/user/project", false},
		{"父目录逃逸", "/home/user/file.txt", "/home/user/project", false},
		{"根目录逃逸", "/etc/passwd", "/home/user/project", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWithinContext(tt.path, tt.contextAbs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCollectFiles_PathTraversal(t *testing.T) {
	// 构建上下文目录
	ctx := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(ctx, "app.go"), []byte("package main"), 0o644))

	// 在上下文之外创建文件
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o644))

	// ../ 相对路径逃逸应被拒绝
	relEscape := "../" + filepath.Base(outside) + "/secret.txt"
	_, err := collectFiles(relEscape, ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "escapes build context")

	// 绝对路径逃逸应被拒绝
	_, err = collectFiles(filepath.Join(outside, "secret.txt"), ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "escapes build context")

	// 合法路径应正常工作
	files, err := collectFiles("app.go", ctx, nil)
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "app.go", files[0].relPath)
}

func TestCollectGlob_FiltersEscapedMatches(t *testing.T) {
	ctx := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(ctx, "a.txt"), []byte("a"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(ctx, "b.txt"), []byte("b"), 0o644))

	// glob 匹配上下文内的文件
	pattern := filepath.Join(ctx, "*.txt")
	files, err := collectGlob(pattern, ctx, nil)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}
