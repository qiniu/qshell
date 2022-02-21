package operations

import (
	"bufio"
	"github.com/qiniu/qshell/v2/iqshell/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"io"
	"os"
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
func Refresh(info RefreshInfo) {
	log.DebugF("qps limit: %d, max item-size: %d", info.QpsLimit, info.SizeLimit)

	var err error
	var urlReader io.ReadCloser
	if len(info.ItemListFile) == 0 {
		urlReader = os.Stdin
	} else {
		urlReader, err = os.Open(info.ItemListFile)
		if err != nil {
			log.ErrorF("Open url list file error:%v", err)
			os.Exit(data.StatusHalt)
		}
	}
	defer urlReader.Close()

	createQpsLimitIfNeeded(info.QpsLimit)

	scanner := bufio.NewScanner(urlReader)
	itemsToRefresh := make([]string, 0, 50)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text())
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
