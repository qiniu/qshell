package operations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/stretchr/testify/require"
)

func sandboxTestAPI(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	t.Setenv("QINIU_API_KEY", "test-key")
	t.Setenv("E2B_API_KEY", "")
	t.Setenv("QINIU_SANDBOX_API_URL", server.URL)
	t.Setenv("SANDBOX_RETRY_MAX", "0")
	data.SetCmdStatus(data.StatusOK)
	t.Cleanup(func() { data.SetCmdStatus(data.StatusOK) })
}

func sandboxOutput(t *testing.T, run func()) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "stdout")
	require.NoError(t, err)
	defer file.Close()
	stdout := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = stdout }()
	run()
	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)
	output, err := io.ReadAll(file)
	require.NoError(t, err)
	return string(output)
}

func TestBulkOperationsEnumerateAllPagesBeforeMutation(t *testing.T) {
	operations := map[string]func(){
		"kill":   func() { Kill(KillInfo{All: true, Metadata: "branch=feat/a+b&x=1%,purpose=page-test"}) },
		"pause":  func() { Pause(PauseInfo{All: true, Metadata: "branch=feat/a+b&x=1%,purpose=page-test"}) },
		"resume": func() { Resume(ResumeInfo{All: true, Metadata: "branch=feat/a+b&x=1%,purpose=page-test"}) },
	}
	for name, run := range operations {
		t.Run(name, func(t *testing.T) {
			var mu sync.Mutex
			pages := 0
			mutatedEarly := false
			mutated := map[string]int{}
			filterChanged := false
			sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes" {
					pages++
					// 服务端先解码过滤串，再解析各个键值。
					raw, err := url.QueryUnescape(r.URL.Query().Get("metadata"))
					filters, parseErr := url.ParseQuery(raw)
					if err != nil || parseErr != nil || len(filters) != 2 ||
						filters.Get("purpose") != "page-test" || filters.Get("branch") != "feat/a+b&x=1%" {
						filterChanged = true
					}
					switch r.URL.Query().Get("nextToken") {
					case "":
						w.Header().Set("X-Next-Token", "page/2?+")
						fmt.Fprint(w, `[{"sandboxID":"sbx-1"}]`)
					case "page/2?+":
						fmt.Fprint(w, `[{"sandboxID":"sbx-2"},{"sandboxID":"sbx-1"}]`)
					default:
						http.Error(w, "unexpected cursor", http.StatusBadRequest)
					}
					return
				}
				if pages != 2 {
					mutatedEarly = true
				}
				path := strings.TrimPrefix(r.URL.Path, "/v2")
				id := strings.Split(strings.TrimPrefix(path, "/sandboxes/"), "/")[0]
				mutated[id]++
				if strings.HasSuffix(r.URL.Path, "/connect") {
					fmt.Fprintf(w, `{"sandboxID":%q,"templateID":"base","clientID":"client","envdVersion":"0.1.0","envdAccessToken":"test-envd-token"}`, id)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			})
			sandboxOutput(t, run)
			require.Equal(t, 2, pages)
			require.False(t, mutatedEarly, "资源变更应在完整枚举后开始")
			require.False(t, filterChanged, "翻页应保留筛选条件")
			require.Equal(t, map[string]int{"sbx-1": 1, "sbx-2": 1}, mutated)
			require.Equal(t, data.StatusOK, data.GetCmdStatus())
		})
	}
}

func TestListAndBulkKillMatchTheSameMetadata(t *testing.T) {
	const metadata = "lang=c++"
	var deleted []string
	var mu sync.Mutex
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			mu.Lock()
			deleted = append(deleted, strings.TrimPrefix(r.URL.Path, "/sandboxes/"))
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
		raw, err := url.QueryUnescape(r.URL.Query().Get("metadata"))
		condition, parseErr := url.ParseQuery(raw)
		if err != nil || parseErr != nil {
			http.Error(w, "invalid metadata", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch condition.Get("lang") {
		case "c++":
			fmt.Fprint(w, `[{"sandboxID":"sbx-plus"}]`)
		case "c  ":
			fmt.Fprint(w, `[{"sandboxID":"sbx-spaces"}]`)
		default:
			fmt.Fprint(w, `[]`)
		}
	})
	for _, all := range []bool{false, true} {
		output := sandboxOutput(t, func() { List(ListInfo{All: all, Metadata: metadata, Format: "json"}) })
		var items []struct{ SandboxID string }
		require.NoError(t, json.Unmarshal([]byte(output), &items))
		require.Len(t, items, 1)
		require.Equal(t, "sbx-plus", items[0].SandboxID)
	}
	sandboxOutput(t, func() { Kill(KillInfo{All: true, Metadata: metadata}) })
	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, []string{"sbx-plus"}, deleted)
	require.Equal(t, data.StatusOK, data.GetCmdStatus())
}

func TestListAndBulkOperationsRejectInvalidMetadataBeforeAnyRequest(t *testing.T) {
	operations := map[string]func(string){
		"list": func(filter string) {
			List(ListInfo{Metadata: filter, MetadataSet: true})
		},
		"kill": func(filter string) {
			Kill(KillInfo{All: true, Metadata: filter, MetadataSet: true})
		},
		"pause": func(filter string) {
			Pause(PauseInfo{All: true, Metadata: filter, MetadataSet: true})
		},
		"resume": func(filter string) {
			Resume(ResumeInfo{All: true, Metadata: filter, MetadataSet: true})
		},
	}
	for name, run := range operations {
		t.Run(name, func(t *testing.T) {
			for _, filter := range []string{"", "env", "env=", "team=ops,env=", "env=prod,typo", "env=prod,env=dev"} {
				t.Run(filter, func(t *testing.T) {
					sandboxTestAPI(t, func(w http.ResponseWriter, _ *http.Request) {
						t.Error("invalid metadata must not send any request")
						http.Error(w, "unexpected request", http.StatusBadRequest)
					})
					sandboxOutput(t, func() { run(filter) })
					require.Equal(t, data.StatusError, data.GetCmdStatus())
				})
			}
		})
	}
}

func TestSandboxManagementFailuresSetExitStatus(t *testing.T) {
	operations := map[string]func(){
		"list":        func() { List(ListInfo{Format: "json"}) },
		"kill":        func() { Kill(KillInfo{SandboxIDs: []string{"sbx-1"}}) },
		"pause":       func() { Pause(PauseInfo{SandboxIDs: []string{"sbx-1"}}) },
		"resume":      func() { Resume(ResumeInfo{SandboxIDs: []string{"sbx-1"}}) },
		"kill-list":   func() { Kill(KillInfo{All: true}) },
		"pause-list":  func() { Pause(PauseInfo{All: true}) },
		"resume-list": func() { Resume(ResumeInfo{All: true}) },
	}
	for name, run := range operations {
		t.Run(name, func(t *testing.T) {
			sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"message":"invalid request"}`)
			})
			sandboxOutput(t, run)
			require.Equal(t, data.StatusError, data.GetCmdStatus())
		})
	}
}

func TestBulkEnumerationFailureDoesNotMutateFirstPage(t *testing.T) {
	var mu sync.Mutex
	mutations := 0
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			mutations++
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.URL.Query().Get("nextToken") == "" {
			w.Header().Set("X-Next-Token", "page-2")
			fmt.Fprint(w, `[{"sandboxID":"sbx-1"}]`)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"second page failed"}`)
	})
	sandboxOutput(t, func() { Kill(KillInfo{All: true}) })
	require.Zero(t, mutations)
	require.Equal(t, data.StatusError, data.GetCmdStatus())
}

func TestListAllPreservesJSONAndPageLimit(t *testing.T) {
	pages := 0
	wrongLimit := false
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		pages++
		wrongLimit = wrongLimit || r.URL.Query().Get("limit") != "1"
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("nextToken") == "" {
			w.Header().Set("X-Next-Token", "page-2")
			fmt.Fprint(w, `[{"sandboxID":"sbx-1"}]`)
			return
		}
		fmt.Fprint(w, `[{"sandboxID":"sbx-2"}]`)
	})
	output := sandboxOutput(t, func() { List(ListInfo{All: true, Limit: 1, Format: "json"}) })
	var items []struct {
		SandboxID string `json:"sandboxID"`
	}
	require.NoError(t, json.Unmarshal([]byte(output), &items))
	require.Equal(t, 2, pages)
	require.False(t, wrongLimit)
	require.Len(t, items, 2)
	require.Equal(t, "sbx-1", items[0].SandboxID)
	require.Equal(t, "sbx-2", items[1].SandboxID)
}

func TestListAcceptsInitialCursor(t *testing.T) {
	cursor := ""
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		cursor = r.URL.Query().Get("nextToken")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})
	sandboxOutput(t, func() { List(ListInfo{NextToken: "page/2?+", Format: "json"}) })
	require.Equal(t, "page/2?+", cursor)
}

func TestListAllRejectsRepeatedCursorWithoutPartialOutput(t *testing.T) {
	pages := 0
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		pages++
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Next-Token", "same-page")
		fmt.Fprint(w, `[{"sandboxID":"sbx-1"}]`)
	})
	output := sandboxOutput(t, func() { List(ListInfo{All: true, Format: "json"}) })
	require.Equal(t, 2, pages)
	require.Empty(t, output)
	require.Equal(t, data.StatusError, data.GetCmdStatus())
}

func TestListReportsCursorOnStderrAndKeepsJSONArray(t *testing.T) {
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Next-Token", "page-2")
		fmt.Fprint(w, `[{"sandboxID":"sbx-1"}]`)
	})
	file, err := os.CreateTemp(t.TempDir(), "stderr")
	require.NoError(t, err)
	defer file.Close()
	stderr := os.Stderr
	os.Stderr = file
	defer func() { os.Stderr = stderr }()
	output := sandboxOutput(t, func() { List(ListInfo{Format: "json"}) })
	var items []map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &items))
	require.Len(t, items, 1)
	message, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Contains(t, string(message), "Next token: page-2")
}

func TestListRejectsAllWithInitialCursor(t *testing.T) {
	requests := 0
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) { requests++ })
	sandboxOutput(t, func() { List(ListInfo{All: true, NextToken: "page-2"}) })
	require.Zero(t, requests)
	require.Equal(t, data.StatusError, data.GetCmdStatus())
}

func TestSandboxManagementRequiresAPIKey(t *testing.T) {
	operations := map[string]func(){
		"list":   func() { List(ListInfo{}) },
		"kill":   func() { Kill(KillInfo{SandboxIDs: []string{"sbx-1"}}) },
		"pause":  func() { Pause(PauseInfo{SandboxIDs: []string{"sbx-1"}}) },
		"resume": func() { Resume(ResumeInfo{SandboxIDs: []string{"sbx-1"}}) },
	}
	for name, run := range operations {
		t.Run(name, func(t *testing.T) {
			sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) { t.Error("unexpected API request") })
			t.Setenv("QINIU_API_KEY", "")
			sandboxOutput(t, run)
			require.Equal(t, data.StatusError, data.GetCmdStatus())
		})
	}
}

func TestBatchPartialFailureContinuesAndReportsCounts(t *testing.T) {
	var mu sync.Mutex
	seen := map[string]int{}
	sandboxTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		id := strings.TrimPrefix(r.URL.Path, "/sandboxes/")
		seen[id]++
		if id == "sbx-2" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"message":"cannot delete"}`)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	output := sandboxOutput(t, func() { Kill(KillInfo{SandboxIDs: []string{"sbx-1", "sbx-2", "sbx-3"}}) })
	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, map[string]int{"sbx-1": 1, "sbx-2": 1, "sbx-3": 1}, seen)
	require.Contains(t, output, "Processed 3 sandboxes: 2 succeeded, 1 failed")
	require.Equal(t, data.StatusError, data.GetCmdStatus())
}
