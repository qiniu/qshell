package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	qshellcmd "github.com/qiniu/qshell/v2/cmd"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBatchStatusCommand 通过本地 API 验证命令参数、SDK 字段传递和导出文件格式。
func TestBatchStatusCommand(t *testing.T) {
	fixtures := map[string]json.RawMessage{
		"full.txt": json.RawMessage(`{"code":200,"data":{
			"fsize":123,"hash":"etag","mimeType":"text/plain","putTime":17000000001234567,"type":5,
			"endUser":"owner","md5":"digest","status":1,"restoreStatus":2,
			"expiration":1706000000,"transitionToIA":1701000000,"transitionToIntelligentTiering":1702000000,
			"transitionToArchiveIR":1703000000,"transitionToARCHIVE":1704000000,"transitionToDeepArchive":1705000000}}`),
		"empty.txt":   json.RawMessage(`{"code":200,"data":{}}`),
		"missing.txt": json.RawMessage(`{"code":612,"data":{"error":"no such file or directory"}}`),
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("X-Reqid", "test-batchstat")
		if r.Method != http.MethodPost || r.URL.Path != "/batch" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse batch request: %v", err)
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		var results []json.RawMessage
		for _, op := range r.PostForm["op"] {
			entry, err := base64.URLEncoding.DecodeString(strings.TrimPrefix(op, "/stat/"))
			key := strings.TrimPrefix(string(entry), "bucket:")
			result, ok := fixtures[key]
			if err != nil || !strings.HasPrefix(string(entry), "bucket:") || !ok {
				t.Errorf("unexpected batch operation: %q", op)
				http.Error(w, "invalid operation", http.StatusBadRequest)
				return
			}
			results = append(results, result)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			t.Errorf("write batch response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	host := strings.TrimPrefix(server.URL, "http://")
	config := fmt.Sprintf(`{"access_key":"test-ak","secret_key":"test-sk","use_https":false,
		"hosts":{"rs":%q,"rsf":%q,"io":%q,"up":%q,"api":%q,"uc":%q},"log":{"log_level":"none"}}`,
		host, host, host, host, host, host)
	require.NoError(t, os.WriteFile(configPath, []byte(config), 0o600))
	executable, err := os.Executable()
	require.NoError(t, err)

	for _, tc := range []struct {
		name       string
		flags      []string
		input      string
		want       string
		wantError  string
		validation bool
	}{
		{
			name: "default six columns", input: "full.txt\n",
			want: "full.txt\t123\tetag\ttext/plain\t17000000001234567\t5\n",
		},
		{
			name: "explicit empty fields", flags: []string{"--show-fields", ""}, input: "full.txt\n",
			want: "full.txt\t123\tetag\ttext/plain\t17000000001234567\t5\n",
		},
		{
			name: "all fields", input: "full.txt\n",
			flags: []string{"--show-fields", "Key,Fsize,Hash,MimeType,PutTime,Type,EndUser,MD5,Status,RestoreStatus,Expiration,TransitionToIA,TransitionToIntelligentTiering,TransitionToArchiveIR,TransitionToArchive,TransitionToDeepArchive"},
			want:  "full.txt\t123\tetag\ttext/plain\t17000000001234567\t5\towner\tdigest\t1\t2\t1706000000\t1701000000\t1702000000\t1703000000\t1704000000\t1705000000\n",
		},
		{
			name: "field order and case", input: "full.txt\n",
			flags: []string{"--show-fields", "expiration,key,md5,Key"}, want: "1706000000\tfull.txt\tdigest\tfull.txt\n",
		},
		{
			name: "leading intermediate and trailing empty columns", input: "empty.txt\n",
			flags: []string{"--show-fields", "EndUser,Key,Status,RestoreStatus,Expiration,TransitionToIA,TransitionToIntelligentTiering,TransitionToArchiveIR,TransitionToArchive,TransitionToDeepArchive,MD5"},
			want:  "\tempty.txt\t\t\t\t\t\t0\t\t\t\n",
		},
		{
			name: "failure isolation", input: "full.txt\nmissing.txt\n",
			flags: []string{"--show-fields", "Key,Expiration"}, want: "full.txt\t1706000000\n", wantError: "no such file or directory",
		},
		{
			name: "unknown field", input: "full.txt\n", flags: []string{"--show-fields", "Key,Unknown"},
			wantError: "show-fields value error:Unknown not support", validation: true,
		},
		{
			name: "empty field", input: "full.txt\n", flags: []string{"--show-fields", "Key,,Expiration"},
			wantError: "show-fields value error: not support", validation: true,
		},
		{
			name: "trailing comma", input: "full.txt\n", flags: []string{"--show-fields", "Key,"},
			wantError: "show-fields value error: not support", validation: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workdir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(workdir, "keys.txt"), []byte(tc.input), 0o600))
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			args := []string{"-test.run=^TestBatchStatusCommandProcess$", "--", "--local", "--config", configPath,
				"batchstat", "bucket", "-i", "keys.txt", "-y", "--worker", "1", "--silence",
				"--outfile", "result.tsv", "--success-list", "success.txt", "--failure-list", "failure.txt"}
			command := exec.CommandContext(ctx, executable, append(args, tc.flags...)...)
			command.Dir = workdir
			command.Env = append(os.Environ(), "QSHELL_BATCHSTAT_COMMAND=1", "QINIU_SDK_DEBUG=false")
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			before := requests.Load()
			err := command.Run()
			require.NoError(t, ctx.Err(), "command timed out")
			if tc.wantError == "" {
				require.NoError(t, err, "stdout: %s\nstderr: %s", &stdout, &stderr)
			} else {
				var exitError *exec.ExitError
				require.ErrorAs(t, err, &exitError)
				assert.Equal(t, data.StatusError, exitError.ExitCode())
				assert.Contains(t, stderr.String(), tc.wantError)
			}
			if tc.validation {
				assert.Equal(t, before, requests.Load(), "invalid fields must not call the API")
				for _, path := range []string{"result.tsv", "success.txt", "failure.txt"} {
					assert.NoFileExists(t, filepath.Join(workdir, path))
				}
				return
			}
			assert.Greater(t, requests.Load(), before)
			result, err := os.ReadFile(filepath.Join(workdir, "result.tsv"))
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(result))
			assert.Contains(t, stdout.String(), tc.want, "console and outfile must use the same format")
			success, err := os.ReadFile(filepath.Join(workdir, "success.txt"))
			require.NoError(t, err)
			assert.Equal(t, strings.ReplaceAll(tc.input, "missing.txt\n", ""), string(success))
			failure, err := os.ReadFile(filepath.Join(workdir, "failure.txt"))
			require.NoError(t, err)
			if tc.wantError != "" {
				assert.Equal(t, "missing.txt\tQShellError:[612]no such file or directory\n", string(failure))
			} else {
				assert.Empty(t, failure)
			}
		})
	}
}

// TestBatchStatusCommandProcess 在子进程执行真实命令，隔离工作区和命令全局状态。
func TestBatchStatusCommandProcess(t *testing.T) {
	if os.Getenv("QSHELL_BATCHSTAT_COMMAND") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"qshell"}, os.Args[i+1:]...)
			qshellcmd.Execute()
			os.Exit(data.GetCmdStatus())
		}
	}
	t.Fatal("missing command arguments")
}
