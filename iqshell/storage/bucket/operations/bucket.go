package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type GetBucketInfo struct {
	Bucket string
}

func (i *GetBucketInfo) Check() *data.CodeError {
	if len(i.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func GetBucket(cfg *iqshell.Config, info GetBucketInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if bucketInfo, err := bucket.GetBucketInfo(bucket.GetBucketApiInfo{
		Bucket: info.Bucket,
	}); err != nil {
		log.ErrorF("get bucket(%s) info error:%v", info.Bucket, err)
	} else {
		desc := ""
		log.AlertF("%-20s:%s", "Bucket", info.Bucket)
		log.AlertF("%-20s:%s", "RegionID", bucketInfo.Region)

		if bucketInfo.Private > 0 {
			desc = "私有空间"
		} else {
			desc = "公有空间"
		}
		log.AlertF("%-20s:%d(%s)", "Private", bucketInfo.Private, desc)

		if bucketInfo.Protected > 0 {
			desc = "原图保护已开启"
		} else {
			desc = "原图保护未开启"
		}
		log.AlertF("%-20s:%d(%s)", "Protected", bucketInfo.Protected, desc)

		if bucketInfo.NoIndexPage > 0 {
			desc = "index.html 不作为默认首页展示"
		} else {
			desc = "index.html 作为默认首页展示"
		}
		log.AlertF("%-20s:%d(%s)", "NoIndexPage", bucketInfo.NoIndexPage, desc)

		log.AlertF("%-20s:%d", "MaxAge", bucketInfo.MaxAge)
		log.AlertF("%-20s:%s", "Separator", bucketInfo.Separator)

		if bucketInfo.TokenAntiLeechMode > 0 {
			desc = "已使用 token 签名的防盗链方式"
		} else {
			desc = "未使用 token 签名的防盗链方式"
		}
		log.AlertF("%-20s:%d(%s)", "TokenAntiLeechMode", bucketInfo.TokenAntiLeechMode, desc)
		log.AlertF("%-20s:%s", "MacKey", bucketInfo.MacKey)
		log.AlertF("%-20s:%s", "MacKey2", bucketInfo.MacKey2)

		if bucketInfo.EnableSource {
			desc = "源站支持防盗链则开启防盗链"
		} else {
			desc = "源站支持防盗链也不开启防盗链"
		}
		log.AlertF("%-20s:%v(%s)", "EnableSource", bucketInfo.EnableSource, desc)

		if bucketInfo.NoRefer {
			desc = "允许空 Referer 访问"
		} else {
			desc = "不允许空 Referer 访问"
		}
		log.AlertF("%-20s:%v(%s)", "NoRefer", bucketInfo.NoRefer, desc)

		if bucketInfo.AntiLeechMode == 0 {
			desc = "未设置防盗链"
		} else if bucketInfo.AntiLeechMode == 1 {
			desc = "设置了防盗链的白名单"
		} else if bucketInfo.AntiLeechMode == 2 {
			desc = "设置了防盗链的黑名单"
		} else {
			desc = ""
		}
		log.AlertF("%-20s:%d(%s)", "AntiLeechMode", bucketInfo.AntiLeechMode, desc)
		log.AlertF("%-20s:", "ReferWhiteList")
		for _, white := range bucketInfo.ReferWl {
			log.AlertF("%-20s%s", "", white)
		}
		log.AlertF("%-20s:", "ReferBlackList")
		for _, black := range bucketInfo.ReferBl {
			log.AlertF("%-20s%s", "", black)
		}

		log.AlertF("%-20s:", "Styles")
		for key, value := range bucketInfo.Styles {
			log.AlertF("%20s%-10s:%s", "", key, value)
		}
	}
}
