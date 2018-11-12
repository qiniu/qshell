package cmd

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api.v7/cdn"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

const (
	BATCH_CDN_REFRESH_URLS_ALLOW_MAX = 100
	BATCH_CDN_REFRESH_DIRS_ALLOW_MAX = 10
	BATCH_CDN_PREFETCH_ALLOW_MAX     = 100
)

var (
	prefetchFile string
	isDir        bool
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
	cdnRefreshCmd.Flags().BoolVarP(&isDir, "dirs", "r", false, "refresh directory")
	cdnRefreshCmd.Flags().StringVarP(&prefetchFile, "input-file", "i", "", "input file")
	cdnPreCmd.Flags().StringVarP(&prefetchFile, "input-file", "i", "", "input file")

	RootCmd.AddCommand(cdnPreCmd, cdnRefreshCmd)
}

func CdnRefresh(cmd *cobra.Command, params []string) {
	var urlListFile string

	if prefetchFile != "" {
		urlListFile = params[0]
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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	cm := iqshell.GetCdnManager()
	scanner := bufio.NewScanner(fp)

	itemsToRefresh := make([]string, 0, 100)

	if isDir {
		for scanner.Scan() {
			item := strings.TrimSpace(scanner.Text())
			if item == "" {
				continue
			}
			itemsToRefresh = append(itemsToRefresh, item)

			if len(itemsToRefresh) == BATCH_CDN_REFRESH_DIRS_ALLOW_MAX {
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

			if len(itemsToRefresh) == BATCH_CDN_REFRESH_URLS_ALLOW_MAX {
				cdnRefresh(cm, itemsToRefresh, nil)
				itemsToRefresh = make([]string, 0, 100)
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
	resp, err := cm.RefreshUrlsAndDirs(urls, dirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CDN refresh error: %v\n", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}

func CdnPrefetch(cmd *cobra.Command, params []string) {
	var urlListFile string

	if prefetchFile != "" {
		urlListFile = params[0]
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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	cm := iqshell.GetCdnManager()
	scanner := bufio.NewScanner(fp)

	urlsToPrefetch := make([]string, 0, 10)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		urlsToPrefetch = append(urlsToPrefetch, url)

		if len(urlsToPrefetch) == BATCH_CDN_PREFETCH_ALLOW_MAX {
			cdnPrefetch(cm, urlsToPrefetch)
			urlsToPrefetch = make([]string, 0, 10)
		}
	}

	if len(urlsToPrefetch) > 0 {
		cdnPrefetch(cm, urlsToPrefetch)
	}
}

func cdnPrefetch(cm *cdn.CdnManager, urls []string) {
	resp, err := cm.PrefetchUrls(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CDN prefetch error: %v\n", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}
