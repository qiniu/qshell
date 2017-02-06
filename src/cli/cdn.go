package cli

import (
	"bufio"
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qshell"
	"strings"
)

const (
	BATCH_CDN_REFRESH_ALLOW_MAX  = 100
	BATCH_CDN_PREFETCH_ALLOW_MAX = 100
)

func CdnRefresh(cmd string, params ...string) {
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

		urlsToRefresh := make([]string, 0, 10)
		for scanner.Scan() {
			url := strings.TrimSpace(scanner.Text())
			if url == "" {
				continue
			}
			urlsToRefresh = append(urlsToRefresh, url)

			if len(urlsToRefresh) == BATCH_CDN_REFRESH_ALLOW_MAX {
				cdnRefresh(&client, urlsToRefresh)
				urlsToRefresh = make([]string, 0, 10)
			}
		}

		if len(urlsToRefresh) > 0 {
			cdnRefresh(&client, urlsToRefresh)
		}
	} else {
		CmdHelp(cmd)
	}
}

func cdnRefresh(client *rs.Client, urls []string) {
	resp, err := qshell.BatchRefresh(client, urls)
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
