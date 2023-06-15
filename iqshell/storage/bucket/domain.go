package bucket

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"strings"
)

func DomainOfBucket(bucket string) (domain string, err *data.CodeError) {
	//get domains of bucket
	domainsOfBucket, gErr := AllDomainsOfBucket(bucket)
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("Get domains of bucket error: %v", gErr)
		return
	}

	if len(domainsOfBucket) == 0 {
		err = data.NewEmptyError().AppendDescF("No domains found for bucket: %s", bucket)
		return
	}
	return domainsOfBucket[0].Domain.Value(), nil
}

var (
	DomainTypeStrings     = []string{"CDN 域名", "源站域名"}
	DomainApiScopeStrings = []string{"kodo api", "s3 api"}
)

type DomainInfo struct {
	Domain      *data.String `json:"domain"`
	DomainType  *data.Int    `json:"domaintype"`   // 0:cdn 1:源站
	ApiScope    *data.Int    `json:"apiscope"`     //
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

func (i DomainInfo) String() string {
	return i.DescriptionString()
}

func (i *DomainInfo) DescriptionString() string {
	return fmt.Sprintf("%s", i.Domain.Value())
}

func (i *DomainInfo) DetailDescriptionString() string {
	return fmt.Sprintf("%-12s: %s\n%-12s: %d(%s)\n%-12s: %d(%s)\n%-12s: %s\n",
		"domain", i.Domain.Value(),
		"type", i.DomainType.Value(), i.getTypeString(),
		"ApiScope", i.ApiScope.Value(), i.getApiScopeString(),
		"FreezeTypes", i.FreezeTypes)
}

// AllDomainsOfBucket 获取一个存储空间绑定的CDN域名
func AllDomainsOfBucket(bucket string) (domains []DomainInfo, err *data.CodeError) {
	domains, err = allDomainsOfBucket(workspace.GetConfig(), bucket)
	if len(domains) == 0 {
		return domains, data.NewEmptyError().AppendDesc("domain list is empty").AppendError(err)
	}

	cdnDomains := make([]DomainInfo, 0)
	sourceDomains := make([]DomainInfo, 0)
	for _, d := range domains {
		if d.Domain != nil && !strings.HasPrefix(d.Domain.Value(), ".") {
			if d.DomainType.Value() == 0 {
				cdnDomains = append(cdnDomains, d)
			} else if d.DomainType.Value() == 1 {
				sourceDomains = append(sourceDomains, d)
			}
		}
	}
	domains = append(cdnDomains, sourceDomains...)
	return domains, err
}

func allDomainsOfBucket(cfg *config.Config, bucket string) ([]DomainInfo, *data.CodeError) {
	bucketManager, gErr := GetBucketManager()
	if gErr != nil {
		return nil, gErr
	}

	var domains []DomainInfo
	reqHost := workspace.GetConfig().Hosts.GetOneUc()
	reqURL := fmt.Sprintf("%s/v3/domains?tbl=%s", utils.Endpoint(cfg.IsUseHttps(), reqHost), bucket)
	//reqURL = fmt.Sprintf("%s/domain?bucket=%s&type=all", utils.Endpoint(cfg.IsUseHttps(), reqHost), bucket)
	err := bucketManager.Client.CredentialedCall(workspace.GetContext(), bucketManager.Mac, auth.TokenQiniu, &domains, "GET", reqURL, nil)
	if err != nil {
		if e, ok := err.(*storage.ErrorInfo); ok {
			if e.Code != 404 {
				return nil, data.NewError(e.Code, "domain list request error:"+e.Err)
			}
		} else {
			return nil, data.NewEmptyError().AppendDesc("domain list request").AppendError(err)
		}
	}
	return domains, nil
}
