package cdn

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

const (
	BatchRefreshUrlsAllowMax = 50 // CDN刷新一次性最大的刷新文件列表
	BatchRefreshDirsAllowMax = 10 // CDN目录刷新一次性最大的刷新目录数
	BatchPrefetchAllowMax    = 50 // 预取一次最大的预取数目
)

// GetCdnManager 获取CdnManager
func getCdnManager() (cdnManager *cdn.CdnManager, err *data.CodeError) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("GetCdnManager error: %v\n", gErr)
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cdnManager = cdn.NewCdnManager(mac)
	return
}

func Prefetch(urls []string) *data.CodeError {
	cdnManager, err := getCdnManager()
	if err != nil {
		return err
	}

	if resp, e := cdnManager.PrefetchUrls(urls); e != nil {
		return data.NewEmptyError().AppendDescF("CDN prefetch error:%v", e)
	} else if resp.Code != 200 {
		return data.NewEmptyError().AppendDescF("CDN prefetch Code: %d, Error: %s", resp.Code, resp.Error)
	} else {
		log.InfoF("CDN prefetch Code: %d, FlowInfo: %s", resp.Code, resp.Error)
	}
	return nil
}

func Refresh(urls []string, dirs []string) *data.CodeError {
	cdnManager, err := getCdnManager()
	if err != nil {
		return err
	}

	log.DebugF("cdnRefresh, url size: %d, dir size: %d", len(urls), len(dirs))
	if resp, e := cdnManager.RefreshUrlsAndDirs(urls, dirs); e != nil {
		return data.NewEmptyError().AppendDescF("CDN refresh error:%v", err)
	} else if resp.Code != 200 {
		return data.NewEmptyError().AppendDescF("CDN refresh Code: %d, Error: %s", resp.Code, resp.Error)
	} else {
		log.InfoF("CDN refresh Code: %d, FlowInfo: %s", resp.Code, resp.Error)
	}
	return nil
}
