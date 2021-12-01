package cdn

import (
	"fmt"
	account2 "github.com/qiniu/qshell/v2/iqshell/account"
	"os"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
)

// 获取CdnManager
func GetCdnManager() *cdn.CdnManager {
	account, gErr := account2.GetAccount()
	if gErr != nil {
		fmt.Fprintf(os.Stderr, "GetCdnManager error: %v\n", gErr)
		os.Exit(1)
	}
	mac := qbox.NewMac(account.AccessKey, account.SecretKey)
	return cdn.NewCdnManager(mac)
}
