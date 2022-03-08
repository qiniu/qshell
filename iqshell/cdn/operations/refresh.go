package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strings"
)

type RefreshInfo struct {
	ItemListFile string
	IsDir        bool
	SizeLimit    int
	QpsLimit     int
}

func (info *RefreshInfo) Check() error {
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

	handler, err := group.NewHandler(group.Info{
		InputFile: info.ItemListFile,
		Force:     true,
	})
	if err != nil {
		log.Error(err)
		return
	}

	createQpsLimitIfNeeded(info.QpsLimit)

	line := ""
	hasMore := false
	itemsToRefresh := make([]string, 0, 50)
	for {
		line, hasMore = handler.Scanner().ScanLine()
		if !hasMore {
			break
		}

		item := strings.TrimSpace(line)
		log.DebugF("read line:%s", item)
		if item == "" {
			continue
		}
		itemsToRefresh = append(itemsToRefresh, item)
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
	var err error

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
	}
	return
}
