package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	cdn2 "github.com/qiniu/qshell/v2/iqshell/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"

	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/spf13/cobra"
)

const (
	// CDN刷新一次性最大的刷新文件列表
	BATCH_CDN_REFRESH_URLS_ALLOW_MAX = 50

	// CDN目录刷新一次性最大的刷新目录数
	BATCH_CDN_REFRESH_DIRS_ALLOW_MAX = 10

	// 预取一次最大的预取数目
	BATCH_CDN_PREFETCH_ALLOW_MAX = 50
)

var (
	prefetchFile string
	isDir        bool
	itemsLimit   int // 每次提交时 url 数量

	timeTicker *time.Ticker
	qpsLimit   int // 每秒 http 请求限制
)

var cdnPreCmd = &cobra.Command{
	Use:   "cdnprefetch [-i <UrlListFile>]",
	Short: "Batch prefetch the urls in the url list file",
	Long:  "Batch prefetch the urls in the url list file or from stdin if UrlListFile not specified",
	Args:  cobra.ExactArgs(0),
	Run:   CdnPrefetch,
}

var cdnRefreshCmd = &cobra.Command{
	Use:   "cdnrefresh [-i <UrlListFile>]",
	Short: "Batch refresh the cdn cache by the url list file",
	Long:  "Batch refresh the cdn cache by the url list file or from stdin if UrlListFile not specified",
	Args:  cobra.ExactArgs(0),
	Run:   CdnRefresh,
}

func init() {
	OnInitialize(initOnInitialize)

	cdnRefreshCmd.Flags().BoolVarP(&isDir, "dirs", "r", false, "refresh directory")
	cdnRefreshCmd.Flags().StringVarP(&prefetchFile, "input-file", "i", "", "input file")
	cdnRefreshCmd.Flags().IntVar(&qpsLimit, "qps", 0, "qps limit for http call")
	cdnRefreshCmd.Flags().IntVarP(&itemsLimit, "size", "s", 0, "max item-size pre commit")

	cdnPreCmd.Flags().StringVarP(&prefetchFile, "input-file", "i", "", "input file")
	cdnPreCmd.Flags().IntVar(&qpsLimit, "qps", 0, "qps limit for http call")
	cdnPreCmd.Flags().IntVarP(&itemsLimit, "size", "s", 0, "max item-size pre commit")

	RootCmd.AddCommand(cdnPreCmd, cdnRefreshCmd)
}

func initOnInitialize() {
	if qpsLimit > 0 {
		d := time.Second / time.Duration(qpsLimit)
		timeTicker = time.NewTicker(d)
	}
}

func acquire() {
	if timeTicker != nil {
		<-timeTicker.C
	}
}

// 【cdnrefresh】刷新所有CDN节点
func CdnRefresh(cmd *cobra.Command, params []string) {
	log.DebugF("qps limit: %d, max item-size: %d", qpsLimit, itemsLimit)

	var urlListFile string

	if prefetchFile != "" {
		urlListFile = prefetchFile
	} else {
		urlListFile = "stdin"
	}

	var fp io.ReadCloser
	var err error

	if urlListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open refresh item list file error,", err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	cm := cdn2.GetCdnManager()
	scanner := bufio.NewScanner(fp)

	itemsToRefresh := make([]string, 0, 100)

	if isDir {
		for scanner.Scan() {
			item := strings.TrimSpace(scanner.Text())
			if item == "" {
				continue
			}
			itemsToRefresh = append(itemsToRefresh, item)

			if len(itemsToRefresh) == BATCH_CDN_REFRESH_DIRS_ALLOW_MAX ||
				(itemsLimit > 0 && len(itemsToRefresh) >= itemsLimit) {
				cdnRefresh(cm, nil, itemsToRefresh)
				itemsToRefresh = make([]string, 0, 10)
			}
		}
	} else {
		for scanner.Scan() {
			item := strings.TrimSpace(scanner.Text())
			if item == "" {
				continue
			}
			itemsToRefresh = append(itemsToRefresh, item)

			if len(itemsToRefresh) == BATCH_CDN_REFRESH_URLS_ALLOW_MAX ||
				(itemsLimit > 0 && len(itemsToRefresh) >= itemsLimit) {
				cdnRefresh(cm, itemsToRefresh, nil)
				itemsToRefresh = make([]string, 0, 50)
			}
		}
	}

	//check final items
	if len(itemsToRefresh) > 0 {
		if isDir {
			cdnRefresh(cm, nil, itemsToRefresh)
		} else {
			cdnRefresh(cm, itemsToRefresh, nil)
		}
	}
}

func cdnRefresh(cm *cdn.CdnManager, urls []string, dirs []string) {
	acquire()
	log.Debug("cdnRefresh, url size: %d, dir size: %d", len(urls), len(dirs))
	resp, err := cm.RefreshUrlsAndDirs(urls, dirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CDN refresh error: %v\n", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}

//  【cdnprefetch】CDN 文件预取
func CdnPrefetch(cmd *cobra.Command, params []string) {
	var urlListFile string

	if prefetchFile != "" {
		urlListFile = prefetchFile
	} else {
		urlListFile = "stdin"
	}

	var fp io.ReadCloser
	var err error

	if urlListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open url list file error,", err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	cm := cdn2.GetCdnManager()
	scanner := bufio.NewScanner(fp)

	urlsToPrefetch := make([]string, 0, 10)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		urlsToPrefetch = append(urlsToPrefetch, url)

		if len(urlsToPrefetch) == BATCH_CDN_PREFETCH_ALLOW_MAX ||
			(itemsLimit > 0 && len(urlsToPrefetch) >= itemsLimit) {
			cdnPrefetch(cm, urlsToPrefetch)
			urlsToPrefetch = make([]string, 0, 10)
		}
	}

	if len(urlsToPrefetch) > 0 {
		cdnPrefetch(cm, urlsToPrefetch)
	}
}

func cdnPrefetch(cm *cdn.CdnManager, urls []string) {
	acquire()
	log.Debug("cdnPrefetch, url size: %d", len(urls))
	resp, err := cm.PrefetchUrls(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CDN prefetch error: %v\n", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}
