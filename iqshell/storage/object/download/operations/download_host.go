package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

func getDownloadDomainAndHost(cfg *config.Config) (domain string, host string) {
	host = getDownloadHost(cfg)
	domain = getDownloadDomain(cfg)
	log.DebugF("get download domain and host, host:%s domain:%s", host, domain)
	if len(host) == 0 && utils.IsIPUrlString(domain) {
		// 如果 domain 配置的为 IP，则查询 bucket 绑定的 host
		if h, err := bucket.DomainOfBucket(cfg.Download.Bucket.Value()); err != nil {
			log.DebugF("get download domain and host, error:%v", err)
		} else {
			host = h
		}
	}
	if len(domain) == 0 {
		domain = host
	}
	return
}

func getDownloadDomain(cfg *config.Config) string {
	// 1. 从 download 配置获取
	domain := cfg.Download.CdnDomain.Value()
	if len(domain) > 0 {
		return domain
	}

	// 2. 从 region 中获取 ioHost
	domain = cfg.Hosts.GetOneIo()
	if len(domain) > 0 {
		return domain
	}

	// 3. 动态获取 bucket 绑定的 domain
	b := cfg.Download.Bucket.Value()
	log.DebugF("get domain of bucket:%s", b)
	if d, e := bucket.DomainOfBucket(b); e != nil {
		log.DebugF("get bucket:%s domain error:%v", b, e)
	} else {
		domain = d
	}
	if len(domain) > 0 {
		return domain
	}

	// 4. 通过 uc query 查询 bucket 所在的 region，并从 region 获取 ioHost
	log.DebugF("get region of bucket:%s", b)
	if region, err := bucket.Region(cfg.Download.Bucket.Value()); err != nil {
		log.DebugF("get region of bucket:%s err:%v", b, err)
	} else {
		domain = region.IovipHost
	}
	return domain
}

func getDownloadHost(cfg *config.Config) string {
	return cfg.Download.IoHost.Value()
}
