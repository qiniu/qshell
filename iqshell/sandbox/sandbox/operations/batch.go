package operations

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// 完整枚举可能包含大量沙箱，限制同时进行的变更请求数。
const maxConcurrentSandboxOperations = 16

func parseMetadataFilter(raw string) (string, error) {
	filters := make(url.Values)
	for pair := range strings.SplitSeq(raw, ",") {
		key, value, ok := strings.Cut(pair, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" || value == "" {
			return "", fmt.Errorf("--metadata requires non-empty key=value pairs")
		}
		key = url.QueryEscape(key)
		if filters.Has(key) {
			return "", fmt.Errorf("--metadata must not contain duplicate keys")
		}
		// 服务端先解码整个过滤串，再分别解码键和值。
		filters.Set(key, url.QueryEscape(value))
	}
	return filters.Encode(), nil
}

func runSandboxBatch(ids []string, verb, pastTense string, operation func(string) error) {
	var succeeded, failed atomic.Int64
	var wg sync.WaitGroup
	jobs := make(chan string)
	for range min(len(ids), maxConcurrentSandboxOperations) {
		wg.Go(func() {
			for id := range jobs {
				if err := operation(id); err != nil {
					failed.Add(1)
					data.SetCmdStatusError()
					sbClient.PrintError("%s sandbox %s failed: %v", verb, id, err)
					continue
				}
				succeeded.Add(1)
				sbClient.PrintSuccess("%s sandbox %s", pastTense, id)
			}
		})
	}
	for _, id := range ids {
		jobs <- id
	}
	close(jobs)
	wg.Wait()
	fmt.Printf("Processed %d sandboxes: %d succeeded, %d failed\n", len(ids), succeeded.Load(), failed.Load())
}
