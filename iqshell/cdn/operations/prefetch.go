package operations

import (
	"bufio"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"os"
	"strings"
)

type PrefetchInfo struct {
	UrlListFile string // url 信息文件
	SizeLimit   int    // 每次刷新最大 size 限制
	QpsLimit    int    // qps 限制
}

func (info *PrefetchInfo) Check() error {
	return nil
}

func Prefetch(cfg *iqshell.Config, info PrefetchInfo) {
	log.DebugF("qps limit: %d, max item-size: %d", info.QpsLimit, info.SizeLimit)

	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	var err error
	var urlReader io.ReadCloser
	if len(info.UrlListFile) == 0 {
		urlReader = os.Stdin
	} else {
		urlReader, err = os.Open(info.UrlListFile)
		if err != nil {
			log.ErrorF("Open url list file error:%v", err)
			os.Exit(data.StatusHalt)
		}
		defer urlReader.Close()
	}

	createQpsLimitIfNeeded(info.QpsLimit)

	scanner := bufio.NewScanner(urlReader)
	urlsToPrefetch := make([]string, 0, 50)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		urlsToPrefetch = append(urlsToPrefetch, url)

		if len(urlsToPrefetch) == cdn.BatchPrefetchAllowMax ||
			(info.SizeLimit > 0 && len(urlsToPrefetch) >= info.SizeLimit) {
			prefetchWithQps(urlsToPrefetch)
			urlsToPrefetch = make([]string, 0, 50)
		}
	}

	if len(urlsToPrefetch) > 0 {
		prefetchWithQps(urlsToPrefetch)
	}
}

func prefetchWithQps(urlsToPrefetch []string) {

	waiterIfNeeded()

	log.Debug("cdnPrefetch, url size: %d", len(urlsToPrefetch))
	if len(urlsToPrefetch) > 0 {
		err := cdn.Prefetch(urlsToPrefetch)
		if err != nil {
			log.Error(err)
		}
	}
}
