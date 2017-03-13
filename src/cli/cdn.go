package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"qshell/qiniu/api.v6/auth/digest"
	"qshell/qiniu/api.v6/rs"
	"strings"
	"qshell/qshell"
)

const (
	BATCH_CDN_REFRESH_URLS_ALLOW_MAX = 100
	BATCH_CDN_REFRESH_DIRS_ALLOW_MAX = 10
	BATCH_CDN_PREFETCH_ALLOW_MAX     = 100
)

func CdnRefresh(cmd string, params ...string) {
	var isDirs bool
	flagSet := flag.NewFlagSet("cdnRefresh", flag.ExitOnError)
	flagSet.BoolVar(&isDirs, "dirs", false, "refresh dirs")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()

	if len(cmdParams) == 1 {
		urlListFile := cmdParams[0]

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
		fp, err := os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open refresh item list file error,", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)

		itemsToRefresh := make([]string, 0, 100)

		if isDirs {
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
			if isDirs {
				cdnRefresh(&client, nil, itemsToRefresh)
			} else {
				cdnRefresh(&client, itemsToRefresh, nil)
			}
		}
	} else {
		CmdHelp(cmd)
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

func CdnPrefetch(cmd string, params ...string) {
	if len(params) == 1 {
		urlListFile := params[0]

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
		fp, err := os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open url list file error,", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
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
	} else {
		CmdHelp(cmd)
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
