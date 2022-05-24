package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strings"
)

type ListBucketInfo struct {
	CToken    string // aws continuation token
	Delimiter string
	MaxKeys   int64
	Prefix    string
	Id        string // id
	SecretKey string
	Region    string
	Bucket    string
}

func (info *ListBucketInfo) Check() *data.CodeError {
	// check AWS bucket
	if info.Bucket == "" {
		return alert.CannotEmptyError("AWS bucket", "")
	}

	// check AWS region
	if info.Region == "" {
		return alert.CannotEmptyError("AWS region", "")
	}

	if info.Id == "" || info.SecretKey == "" {
		return alert.CannotEmptyError("AWS ID and SecretKey", "")
	}

	if info.MaxKeys <= 0 || info.MaxKeys > 1000 {
		log.WarningF("max key:%d out of range [0, 1000], change to 1000", info.MaxKeys)
		info.MaxKeys = 1000
	}

	return nil
}

func ListBucket(cfg *iqshell.Config, info ListBucketInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	log.Alert("Key\tSize\tETag\tLastModified")
	if err := listBucket(info, func(s3 *s3.S3, object *s3.Object) {
		log.AlertF("%s\t%d\t%s\t%s", *object.Key, *object.Size, *object.ETag, *object.LastModified)
	}); err != nil {
		log.Error(err)
	}
}

func listBucket(info ListBucketInfo, objectHandler func(s3 *s3.S3, object *s3.Object)) *data.CodeError {
	if objectHandler == nil {
		return nil
	}

	// AWS related code
	s3session, err := session.NewSession()
	if err != nil {
		return data.NewEmptyError().AppendDescF("create AWS session error:%v", err)
	}
	s3session.Config.WithRegion(info.Region)
	s3session.Config.WithCredentials(credentials.NewStaticCredentials(info.Id, info.SecretKey, ""))

	svc := s3.New(s3session)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(info.Bucket),
		Prefix:  aws.String(info.Prefix),
		MaxKeys: aws.Int64(info.MaxKeys),
	}

	if info.CToken != "" {
		input.ContinuationToken = aws.String(info.CToken)
	}

	for {
		result, lErr := svc.ListObjectsV2(input)

		if lErr != nil {
			if aErr, ok := lErr.(awserr.Error); ok {
				switch aErr.Code() {
				case s3.ErrCodeNoSuchBucket:
					return data.NewEmptyError().AppendDescF("list error:%s %v ContinuationToken:", s3.ErrCodeNoSuchBucket, aErr.Error(), input.ContinuationToken)
				default:
					return data.NewEmptyError().AppendDescF("list error:%v ContinuationToken:", aErr.Error(), input.ContinuationToken)
				}
			} else {
				return data.NewEmptyError().AppendDescF("list error:%v ContinuationToken:", aErr.Error(), input.ContinuationToken)
			}
		}

		for _, obj := range result.Contents {
			if strings.HasSuffix(*obj.Key, "/") && *obj.Size == 0 { // 跳过目录
				continue
			}
			objectHandler(svc, obj)
		}

		if *result.IsTruncated {
			input.ContinuationToken = result.NextContinuationToken
		} else {
			break
		}
	}
	return nil
}
