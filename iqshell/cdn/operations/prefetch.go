package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strings"
)

type PrefetchInfo struct {
	UrlListFile string // url 信息文件
	SizeLimit   int    // 每次刷新最大 size 限制
	QpsLimit    int    // qps 限制
}

func (info *PrefetchInfo) Check() *data.CodeError {
	return nil
}

func Prefetch(cfg *iqshell.Config, info PrefetchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	log.DebugF("qps limit: %d, max item-size: %d", info.QpsLimit, info.SizeLimit)

	workProvider, err := flow.NewWorkProviderOfFile(info.UrlListFile,
		true,
		flow.NewItemsWorkCreator(flow.DefaultLineItemSeparate,
			1,
			func(items []string) (work flow.Work, err *data.CodeError) {
				item := strings.TrimSpace(items[0])
				if item == "" {
					return nil, alert.Error("url invalid", "")
				}
				return &prefetchWork{
					Url: item,
				}, nil
			}))
	if err != nil {
		log.Error(err)
		return
	}

	createQpsLimitIfNeeded(info.QpsLimit)

	urlsToPrefetch := make([]string, 0, 50)
	for {
		hasMore, workInfo, pErr := workProvider.Provide()
		if workInfo == nil || workInfo.Work == nil || pErr != nil {
			log.ErrorF("read work error:%v", pErr)
			continue
		}
		if !hasMore {
			break
		}

		w, _ := workInfo.Work.(*prefetchWork)
		urlsToPrefetch = append(urlsToPrefetch, w.Url)

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

type prefetchWork struct {
	Url string
}

func (w *prefetchWork) WorkId() string {
	return w.Url
}
