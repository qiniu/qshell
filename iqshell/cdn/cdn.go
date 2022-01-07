package cdn

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

const (
	BatchRefreshUrlsAllowMax = 50 // CDN刷新一次性最大的刷新文件列表
	BatchRefreshDirsAllowMax = 10 // CDN目录刷新一次性最大的刷新目录数
	BatchPrefetchAllowMax    = 50 // 预取一次最大的预取数目
)

// GetCdnManager 获取CdnManager
func getCdnManager() (cdnManager *cdn.CdnManager, err error) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = errors.New(fmt.Sprintf("GetCdnManager error: %v\n", gErr))
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cdnManager = cdn.NewCdnManager(mac)
	return
}

func Prefetch(urls []string) (err error) {
	cdnManager, err := getCdnManager()
	if err != nil {
		return err
	}

	resp, err := cdnManager.PrefetchUrls(urls)
	if err != nil {
		err = errors.New("CDN prefetch error:" + err.Error())
	} else if resp.Code != 200 {
		err = errors.New(fmt.Sprintf("CDN prefetch Code: %d, Error: %s", resp.Code, resp.Error))
	} else {
		log.InfoF("CDN prefetch Code: %d, Info: %s", resp.Code, resp.Error)
	}
	return
}

func Refresh(urls []string, dirs []string) (err error) {
	cdnManager, err := getCdnManager()
	if err != nil {
		return err
	}

	log.DebugF("cdnRefresh, url size: %d, dir size: %d", len(urls), len(dirs))
	resp, err := cdnManager.RefreshUrlsAndDirs(urls, dirs)
	if err != nil {
		err = errors.New("CDN refresh error:" + err.Error())
	} else if resp.Code != 200 {
		err = errors.New(fmt.Sprintf("CDN refresh Code: %d, Error: %s", resp.Code, resp.Error))
	} else {
		log.InfoF("CDN refresh Code: %d, Info: %s", resp.Code, resp.Error)
	}
	return
}
