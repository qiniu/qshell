package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"github.com/tonycai653/iqshell/qshell"
	"io"
	"os"
	"strings"
)

const (
	BATCH_CDN_REFRESH_URLS_ALLOW_MAX = 100
	BATCH_CDN_REFRESH_DIRS_ALLOW_MAX = 10
	BATCH_CDN_PREFETCH_ALLOW_MAX     = 100
)

var cdnPreCmd = &cobra.Command{
	Use:   "cdnprefetch [<UrlListFile>]",
	Short: "Batch prefetch the urls in the url list file",
	Long:  "Batch prefetch the urls in the url list file or from stdin if UrlListFile not specified",
	Args:  cobra.RangeArgs(0, 1),
	Run:   CdnPrefetch,
}

var cdnRefreshCmd = &cobra.Command{
	Use:   "cdnrefresh [<UrlListFile>]",
	Short: "Batch refresh the cdn cache by the url list file",
	Long:  "Batch refresh the cdn cache by the url list file or from stdin if UrlListFile not specified",
	Args:  cobra.RangeArgs(0, 1),
	Run:   CdnRefresh,
}

var (
	isDir bool
)

func init() {
	cdnRefreshCmd.Flags().BoolVarP(&isDir, "dirs", "r", false, "refresh directory")

	RootCmd.AddCommand(cdnPreCmd, cdnRefreshCmd)
}

func CdnRefresh(cmd *cobra.Command, params []string) {
	var urlListFile string

	if len(params) == 1 {
		urlListFile = params[0]
	} else {
		urlListFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}

	client := rs.NewMac(&mac)

	var fp io.ReadCloser
	var err error

	if urlListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open refresh item list file error,", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	itemsToRefresh := make([]string, 0, 100)

	if isDir {
		for scanner.Scan() {
			item := strings.TrimSpace(scanner.Text())
			if item == "" {
				continue
			}
			itemsToRefresh = append(itemsToRefresh, item)

			if len(itemsToRefresh) == BATCH_CDN_REFRESH_DIRS_ALLOW_MAX {
				cdnRefresh(&client, nil, itemsToRefresh)
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
				cdnRefresh(&client, itemsToRefresh, nil)
				itemsToRefresh = make([]string, 0, 100)
			}
		}
	}

	//check final items
	if len(itemsToRefresh) > 0 {
		if isDir {
			cdnRefresh(&client, nil, itemsToRefresh)
		} else {
			cdnRefresh(&client, itemsToRefresh, nil)
		}
	}
}

func cdnRefresh(client *rs.Client, urls []string, dirs []string) {
	resp, err := qshell.BatchRefresh(client, urls, dirs)
	if err != nil {
		fmt.Println("CDN refresh error,", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}

func CdnPrefetch(cmd *cobra.Command, params []string) {
	var urlListFile string

	if len(params) == 1 {
		urlListFile = params[0]
	} else {
		urlListFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}

	client := rs.NewMac(&mac)

	var fp io.ReadCloser
	var err error

	if urlListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open url list file error,", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	urlsToPrefetch := make([]string, 0, 10)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		urlsToPrefetch = append(urlsToPrefetch, url)

		if len(urlsToPrefetch) == BATCH_CDN_PREFETCH_ALLOW_MAX {
			cdnPrefetch(&client, urlsToPrefetch)
			urlsToPrefetch = make([]string, 0, 10)
		}
	}

	if len(urlsToPrefetch) > 0 {
		cdnPrefetch(&client, urlsToPrefetch)
	}
}

func cdnPrefetch(client *rs.Client, urls []string) {
	resp, err := qshell.BatchPrefetch(client, urls)
	if err != nil {
		fmt.Println("CDN prefetch error,", err)
	} else {
		if resp.Error != "" {
			fmt.Println(fmt.Sprintf("Code: %d, Info: %s", resp.Code, resp.Error))
		}
	}
}
