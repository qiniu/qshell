package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

func getDownloadDomainAndHost(cfg *config.Config, downloadCfg *DownloadCfg) (domain string, host string) {
	domain = getDownloadDomain(cfg, downloadCfg)
	log.DebugF("get download domain and host, host:%s domain:%s", host, domain)
	if utils.IsIPUrlString(domain) {
		// 如果 domain 配置的为 IP，则查询 bucket 绑定的 host
		if h, err := bucket.DomainOfBucket(downloadCfg.Bucket); err != nil {
			log.DebugF("get download domain and host, error:%v", err)
		} else {
			host = h
		}
	}
	return
}

func getDownloadDomain(cfg *config.Config, downloadCfg *DownloadCfg) string {
	if downloadCfg.GetFileApi {
		return getFileApiDomain(cfg, downloadCfg)
	} else {
		return defaultDownloadDomain(cfg, downloadCfg)
	}
}

func defaultDownloadDomain(cfg *config.Config, downloadCfg *DownloadCfg) string {
	// 1. 从 download 配置获取
	domain := downloadCfg.CdnDomain
	if len(domain) > 0 {
		return domain
	}

	// 2. 动态获取 bucket 绑定的 domain
	b := downloadCfg.Bucket
	log.DebugF("get domain of bucket:%s", b)
	if d, e := bucket.DomainOfBucket(b); e != nil {
		log.DebugF("get bucket:%s domain error:%v", b, e)
	} else {
		domain = d
	}
	if len(domain) > 0 {
		return domain
	}

	return domain
}

func getFileApiDomain(cfg *config.Config, downloadCfg *DownloadCfg) string {
	// 1. 从 download 配置获取
	domain := downloadCfg.CdnDomain
	if len(domain) > 0 {
		return domain
	}

	domain = downloadCfg.IoHost
	if len(domain) > 0 {
		return domain
	}

	// 2. 从 region 中获取 ioHost
	domain = cfg.Hosts.GetOneIo()
	if len(domain) > 0 {
		return domain
	}

	// 3. 通过 uc query 查询 bucket 所在的 region，并从 region 获取 ioHost
	log.DebugF("get region of bucket:%s", downloadCfg.Bucket)
	if region, err := bucket.Region(downloadCfg.Bucket); err != nil {
		log.DebugF("get region of bucket:%s err:%v", downloadCfg.Bucket, err)
	} else {
		domain = region.IovipHost
	}

	return domain
}
