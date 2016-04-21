package cli

import (
	"fmt"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qiniu/log"
	"qshell"
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
			fmt.Println(cdnInfo.LineCname)
			fmt.Println(cdnInfo.CdnInfo)
			fmt.Println(cdnInfo.IpAddress)
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
