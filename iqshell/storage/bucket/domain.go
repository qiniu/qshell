package bucket

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"strings"
)

func DomainOfBucket(bucket string) (domain string, err error) {
	//get domains of bucket
	domainsOfBucket, gErr := AllDomainsOfBucket(bucket)
	if gErr != nil {
		err = fmt.Errorf("Get domains of bucket error: %v", gErr)
		return
	}

	if len(domainsOfBucket) == 0 {
		err = fmt.Errorf("No domains found for bucket: %s", bucket)
		return
	}

	for _, d := range domainsOfBucket {
		if d.Domain != nil && !strings.HasPrefix(d.Domain.Value(), ".") {
			domain = d.Domain.Value()
			break
		}
	}
	return
}

var (
	DomainTypeStrings     = []string{"CDN 域名", "源站域名"}
	DomainApiScopeStrings = []string{"kodo api", "s3 api"}
)

type DomainInfo struct {
	Domain      *data.String `json:"domain"`
	DomainType  *data.Int    `json:"domaintype"`
	ApiScope    *data.Int    `json:"apiscope"`
	FreezeTypes []string     `json:"freeze_types"` // 不为空表示已被冻结
	Tbl         *data.String `json:"tbl"`          // 存储空间名字
	Owner       *data.Int    `json:"uid"`          // 用户UID
	Refresh     *data.Bool   `json:"refresh"`      // cdn的自主刷新
	Ctime       *data.Int    `json:"ctime"`
	Utime       *data.Int    `json:"utime"`
}

func (i *DomainInfo) getTypeString() string {
	if i.DomainType.Value() < 0 || i.DomainType.Value() > len(DomainTypeStrings) {
		return "Unknown"
	}
	return DomainTypeStrings[i.DomainType.Value()]
}

func (i *DomainInfo) getApiScopeString() string {
	if i.ApiScope.Value() < 0 || i.ApiScope.Value() > len(DomainTypeStrings) {
		return "Unknown"
	}
	return DomainApiScopeStrings[i.ApiScope.Value()]
}

func (i *DomainInfo) DescriptionString() string {
	return fmt.Sprintf("%s", i.Domain.Value())
}

func (i *DomainInfo) DetailDescriptionString() string {
	i.FreezeTypes = []string{"1", "2"}
	return fmt.Sprintf("%s\n%-12s: %s(%d)\n%-12s: %s(%d)\n%-12s: %s\n",
		i.Domain.Value(),
		"type", i.getTypeString(), i.DomainType.Value(),
		"ApiScope", i.getApiScopeString(), i.ApiScope.Value(),
		"FreezeTypes", i.FreezeTypes)
}

// AllDomainsOfBucket 获取一个存储空间绑定的CDN域名
func AllDomainsOfBucket(bucket string) (domains []DomainInfo, err error) {
	return allDomainsOfBucket(workspace.GetConfig(), bucket)
}

func allDomainsOfBucket(cfg *config.Config, bucket string) (domains []DomainInfo, err error) {
	bucketManager, gErr := GetBucketManager()
	if gErr != nil {
		return nil, gErr
	}

	reqHost := cfg.Hosts.GetOneApi()
	if len(reqHost) == 0 {
		if region, err := Region(bucket); err != nil {
			return nil, alert.Error("get region error:"+err.Error(), "")
		} else {
			reqHost = region.ApiHost
		}
	}

	reqURL := fmt.Sprintf("%s/v7/domain/list?tbl=%s", utils.Endpoint(cfg.IsUseHttps(), reqHost), bucket)
	err = bucketManager.Client.CredentialedCall(workspace.GetContext(), bucketManager.Mac, auth.TokenQiniu, &domains, "GET", reqURL, nil)
	if err != nil {
		if e, ok := err.(*storage.ErrorInfo); ok {
			if e.Code != 404 {
				return
			}
			err = nil
		} else {
			return
		}
	}
	return domains, nil
}
