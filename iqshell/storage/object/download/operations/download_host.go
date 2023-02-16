package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/host"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"strings"
)

func getDownloadHostProvider(cfg *config.Config, downloadCfg *DownloadCfg) host.Provider {
	hosts := getDownloadHosts(cfg, downloadCfg)
	return host.NewListProvider(hosts)
}

func getDownloadHosts(cfg *config.Config, downloadCfg *DownloadCfg) []*host.Host {
	var hosts []*host.Host
	if downloadCfg.GetFileApi {
		hosts = getFileApiHosts(cfg, downloadCfg)
	} else {
		hosts = defaultDownloadHosts(cfg, downloadCfg)
	}

	hostStrings := make([]string, 0, len(hosts))
	for _, h := range hosts {
		hostStrings = append(hostStrings, h.GetServer())
	}
	log.DebugF("download Domain:[%s]", strings.Join(hostStrings, ","))
	return hosts
}

func defaultDownloadHosts(cfg *config.Config, downloadCfg *DownloadCfg) []*host.Host {

	hosts := make([]*host.Host, 0)

	// 1. 从 download 配置获取
	if len(downloadCfg.CdnDomain) > 0 {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: downloadCfg.CdnDomain,
		})
	}

	// 2. 动态获取 bucket 绑定的 domain
	b := downloadCfg.Bucket
	if domains, e := bucket.AllDomainsOfBucket(b); e != nil {
		log.DebugF("get bucket:%s domain error:%v", b, e)
	} else {
		log.DebugF("get domain of bucket:%s domains:%s", b, domains)
		for _, domain := range domains {
			if data.NotEmpty(domain.Domain) {
				hosts = append(hosts, &host.Host{
					Host:   "",
					Domain: domain.Domain.Value(),
				})
			}
		}
	}

	// 3. 源站域名
	// 此处不传 RegionId 通过 UC bucket 接口自动获取
	if domain, e := download.CreateSrcDownloadDomainWithBucket(cfg, downloadCfg.Bucket, ""); e != nil {
		log.DebugF("create bucket:%s src domain error:%v", b, e)
	} else {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: domain,
		})
	}

	return hosts
}

func getFileApiHosts(cfg *config.Config, downloadCfg *DownloadCfg) []*host.Host {
	hosts := make([]*host.Host, 0)

	// 1. 从 download 配置获取
	if len(downloadCfg.CdnDomain) > 0 {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: downloadCfg.CdnDomain,
		})
	}

	if len(downloadCfg.IoHost) > 0 {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: downloadCfg.IoHost,
		})
	}

	// 2. 从 region 中获取 ioHost
	for _, io := range cfg.Hosts.Io {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: io,
		})
	}

	// 3. 通过 uc query 查询 bucket 所在的 region，并从 region 获取 ioHost
	log.DebugF("get region of bucket:%s", downloadCfg.Bucket)
	if region, err := bucket.Region(downloadCfg.Bucket); err != nil {
		log.DebugF("get region of bucket:%s err:%v", downloadCfg.Bucket, err)
	} else {
		hosts = append(hosts, &host.Host{
			Host:   "",
			Domain: region.IovipHost,
		})
	}

	return hosts
}
