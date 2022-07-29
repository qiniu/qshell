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

type RefreshInfo struct {
	ItemListFile string
	IsDir        bool
	SizeLimit    int
	QpsLimit     int
}

func (info *RefreshInfo) Check() *data.CodeError {
	return nil
}

// Refresh 【cdnrefresh】刷新所有CDN节点
func Refresh(cfg *iqshell.Config, info RefreshInfo) {

	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	log.DebugF("qps limit: %d, max item-size: %d", info.QpsLimit, info.SizeLimit)

	workProvider, err := flow.NewWorkProviderOfFile(info.ItemListFile,
		true,
		flow.NewItemsWorkCreator(flow.DefaultLineItemSeparate,
			1,
			func(items []string) (work flow.Work, err *data.CodeError) {
				item := strings.TrimSpace(items[0])
				if item == "" {
					return nil, alert.Error("url invalid", "")
				}
				return &refreshWork{
					Url: item,
				}, nil
			}))
	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
		return
	}

	createQpsLimitIfNeeded(info.QpsLimit)

	itemsToRefresh := make([]string, 0, 50)
	for {
		hasMore, workInfo, pErr := workProvider.Provide()
		if pErr != nil {
			data.SetCmdStatusError()
			log.ErrorF("read work error:%v", pErr)
			continue
		}
		if !hasMore {
			break
		}

		w, _ := workInfo.Work.(*refreshWork)
		itemsToRefresh = append(itemsToRefresh, w.Url)
		if refreshWithQps(info, itemsToRefresh, false) {
			itemsToRefresh = make([]string, 0, 50)
		}
	}

	//check final items
	if len(itemsToRefresh) > 0 {
		refreshWithQps(info, itemsToRefresh, true)
	}
}

func refreshWithQps(info RefreshInfo, items []string, force bool) (isRefresh bool) {
	var err *data.CodeError

	if info.IsDir {
		if force || len(items) == cdn.BatchRefreshDirsAllowMax ||
			(info.SizeLimit > 0 && len(items) >= info.SizeLimit) {
			waiterIfNeeded()
			err = cdn.Refresh(nil, items)
			isRefresh = true
		}
	} else {
		if force || len(items) == cdn.BatchRefreshUrlsAllowMax ||
			(info.SizeLimit > 0 && len(items) >= info.SizeLimit) {
			waiterIfNeeded()
			err = cdn.Refresh(items, nil)
			isRefresh = true
		}
	}

	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
	}
	return
}

type refreshWork struct {
	Url string
}

func (w *refreshWork) WorkId() string {
	return w.Url
}
