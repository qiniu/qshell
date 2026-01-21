package bucket

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

type CreateApiInfo struct {
	RegionId string
	Bucket   string
	Private  bool
}

func Create(info CreateApiInfo) *data.CodeError {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return err
	}

	cfg := workspace.GetConfig()
	ucHost := cfg.Hosts.GetOneUc()
	if len(ucHost) == 0 {
		return data.NewEmptyError().AppendDesc("can't get uc host")
	}

	url := utils.Endpoint(cfg.IsUseHttps(), ucHost)
	reqURL := fmt.Sprintf("%s/mkbucketv3/%s/region/%s/private/%v", url, info.Bucket, info.RegionId, info.Private)
	e := bucketManager.Client.CredentialedCall(context.Background(), bucketManager.Mac, auth.TokenQiniu, nil, "POST", reqURL, nil)
	return data.ConvertError(e)
}
