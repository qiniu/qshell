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
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"strings"
	"sync"
)

type FetchInfo struct {
	QiniuBucket   string
	Host          string
	BatchInfo     group.Info
	AwsBucketInfo ListBucketInfo
}

func (info *FetchInfo) Check() error {
	if info.AwsBucketInfo.Id == "" || info.AwsBucketInfo.SecretKey == "" {
		return alert.CannotEmptyError("AWS ID and SecretKey", "")
	}

	if info.AwsBucketInfo.Region == "" {
		return alert.CannotEmptyError("AWS region", "")
	}

	if info.AwsBucketInfo.MaxKeys <= 0 || info.AwsBucketInfo.MaxKeys > 1000 {
		log.Warning("max key:%d out of range {}, change to 1000", info.AwsBucketInfo.MaxKeys)
		info.AwsBucketInfo.MaxKeys = 1000
	}

	if info.BatchInfo.WorkCount <= 0 || info.BatchInfo.WorkCount >= 1000 {
		info.BatchInfo.WorkCount = 20
	}

	return nil
}

func Fetch(info FetchInfo) {
	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:     info.BatchInfo.FailExportFilePath,
		OverrideExportFilePath: info.BatchInfo.OverrideExportFilePath,
	})

	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	fetchInfoChan := make(chan object.FetchApiInfo, info.BatchInfo.WorkCount)
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
				fetchInfoChan <- object.FetchApiInfo{
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
	waiter.Add(info.BatchInfo.WorkCount)
	for i := 0; i < info.BatchInfo.WorkCount; i++ {
		go func() {
			for info := range fetchInfoChan {
				_, err := object.Fetch(info)
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
