package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs/operations"
	"strings"
	"sync"
)

type FetchInfo struct {
	QiniuBucket   string
	Host          string
	BatchInfo     operations.BatchInfo
	AwsBucketInfo ListBucketInfo
}

func Fetch(info FetchInfo) {
	if info.AwsBucketInfo.Id == "" || info.AwsBucketInfo.SecretKey == "" {
		log.Error(alert.CannotEmpty("AWS ID and SecretKey", ""))
		return
	}

	if info.AwsBucketInfo.MaxKeys <= 0 || info.AwsBucketInfo.MaxKeys > 1000 {
		log.Warning("max key:%d out of range {}, change to 1000", info.AwsBucketInfo.MaxKeys)
		info.AwsBucketInfo.MaxKeys = 1000
	}

	// check AWS region
	if info.AwsBucketInfo.Region == "" {
		log.Error(alert.CannotEmpty("AWS region", ""))
		return
	}

	if info.BatchInfo.Worker <= 0 || info.BatchInfo.Worker >= 1000 {
		info.BatchInfo.Worker = 20
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:     info.BatchInfo.FailExportFilePath,
		OverrideExportFilePath: info.BatchInfo.OverrideExportFilePath,
	})

	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	fetchInfoChan := make(chan rs.FetchApiInfo, info.BatchInfo.Worker)
	// 生产者
	go func() {
		// AWS related code
		s3session := session.New()
		s3session.Config.WithRegion(info.AwsBucketInfo.Bucket)
		s3session.Config.WithCredentials(credentials.NewStaticCredentials(info.AwsBucketInfo.Id, info.AwsBucketInfo.SecretKey, ""))

		svc := s3.New(s3session)
		input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(info.AwsBucketInfo.Bucket),
			Prefix:    aws.String(info.AwsBucketInfo.Prefix),
			Delimiter: aws.String(info.AwsBucketInfo.Delimiter),
			MaxKeys:   aws.Int64(info.AwsBucketInfo.MaxKeys),
		}
		if info.AwsBucketInfo.CToken != "" {
			input.ContinuationToken = aws.String(info.AwsBucketInfo.CToken)
		}

		for {
			result, err := svc.ListObjectsV2(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case s3.ErrCodeNoSuchBucket:
						log.ErrorF(s3.ErrCodeNoSuchBucket, aerr.Error())
					default:
						log.Error(aerr.Error())
					}
				} else {
					log.Error(err.Error())
				}
				log.ErrorF("ContinuationToken: %v", input.ContinuationToken)
				break
			}

			for _, obj := range result.Contents {
				if strings.HasSuffix(*obj.Key, "/") && *obj.Size == 0 { // 跳过目录
					continue
				}
				fetchInfoChan <- rs.FetchApiInfo{
					Bucket:  info.QiniuBucket,
					Key:     *obj.Key,
					FromUrl: awsUrl(info.AwsBucketInfo.Bucket, info.AwsBucketInfo.Region, *obj.Key),
				}
			}

			if *result.IsTruncated {
				input.ContinuationToken = result.NextContinuationToken
			} else {
				break
			}
		}
		close(fetchInfoChan)
	}()

	// 消费者
	waiter := sync.WaitGroup{}
	waiter.Add(info.BatchInfo.Worker)
	for i := 0; i < info.BatchInfo.Worker; i++ {
		go func() {
			for info := range fetchInfoChan {
				_, err := rs.Fetch(info)
				if err != nil {
					log.ErrorF("fetch %s => %s:%s failed", info.FromUrl, info.Bucket, info.Key)
					resultExport.Fail().ExportF("%s\t%s\t%v", info.FromUrl, info.Key, err)
				} else {
					log.AlertF("fetch %s => %s:%s success", info.FromUrl, info.Bucket, info.Key)
					resultExport.Success().ExportF("%s\t%s", info.FromUrl, info.Key)
				}
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
}

func awsUrl(awsBucket, region, key string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", awsBucket, region, key)
}
