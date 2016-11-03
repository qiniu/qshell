package cli

import (
	"bufio"
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qiniu/log"
	"qshell"
	"strings"
)

const (
	BATCH_CDN_REFRESH_ALLOW_MAX  = 100
	BATCH_CDN_PREFETCH_ALLOW_MAX = 100
)

func GetCdnSupplierOfIp(cmd string, params ...string) {
	if len(params) == 1 {
		ip := params[0]
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		rsClient := rs.NewMac(&mac)
		cdnInfo, err := qshell.GetCdnSupplierOfIp(&rsClient, ip)
		if err != nil {
			log.Error("Get cdn supplier of ip error,", err)
		} else {
			//for _, cname := range cdnInfo.LineCname {
			//	fmt.Println(cname)
			//}
			fmt.Println(cdnInfo.CdnInfo)
			//fmt.Println(cdnInfo.IpAddress)
		}
	} else {
		CmdHelp(cmd)
	}
}

func GetCdnRegionalIps(cmd string, params ...string) {
	if len(params) == 3 {
		cname := params[0]
		isp := params[1]
		province := params[2]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		rsClient := rs.NewMac(&mac)
		ips, err := qshell.GetCdnRegionalIps(&rsClient, cname, isp, province)
		if err != nil {
			log.Error("Get regional ips of cname error,", err)
		} else {
			for _, ip := range ips {
				fmt.Println(ip)
			}
		}
	} else {
		CmdHelp(cmd)
	}
}

func CdnRefresh(cmd string, params ...string) {
	if len(params) == 1 {
		urlListFile := params[0]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}

		client := rs.NewMac(&mac)
		fp, err := os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open url list file error,", err)
			return
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

		fmt.Println("All refresh requests sent!")
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}

		client := rs.NewMac(&mac)
		fp, err := os.Open(urlListFile)
		if err != nil {
			fmt.Println("Open url list file error,", err)
			return
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

		fmt.Println("All prefetch requests sent!")
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
