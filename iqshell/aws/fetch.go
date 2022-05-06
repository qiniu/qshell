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
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strings"
	"sync"
	"time"
)

type FetchInfo struct {
	QiniuBucket   string
	Host          string
	BatchInfo     batch.Info
	AwsBucketInfo ListBucketInfo
}

func (info *FetchInfo) Check() *data.CodeError {
	// check AWS bucket
	if info.AwsBucketInfo.Bucket == "" {
		return alert.CannotEmptyError("AWS bucket", "")
	}

	// check AWS region
	if info.AwsBucketInfo.Region == "" {
		return alert.CannotEmptyError("AWS region", "")
	}

	// check AWS region
	if info.QiniuBucket == "" {
		return alert.CannotEmptyError("Qiniu bucket", "")
	}

	if info.AwsBucketInfo.Id == "" || info.AwsBucketInfo.SecretKey == "" {
		return alert.CannotEmptyError("AWS ID and SecretKey", "")
	}

	if info.BatchInfo.WorkerCount <= 0 || info.BatchInfo.WorkerCount >= 1000 {
		info.BatchInfo.WorkerCount = 20
	}

	return nil
}

func Fetch(cfg *iqshell.Config, info FetchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:   info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:      info.BatchInfo.FailExportFilePath,
		OverwriteExportFilePath: info.BatchInfo.OverwriteExportFilePath,
	})

	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	fetchInfoChan := make(chan object.FetchApiInfo, info.BatchInfo.WorkerCount)
	// 生产者
	go func() {
		// AWS related code
		s3session, nErr := session.NewSession()
		if nErr != nil {
			log.ErrorF("create AWS session error:%v", nErr)
			return
		}
		s3session.Config.WithRegion(info.AwsBucketInfo.Region)
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
			result, lErr := svc.ListObjectsV2(input)
			if lErr != nil {
				if aErr, ok := lErr.(awserr.Error); ok {
					switch aErr.Code() {
					case s3.ErrCodeNoSuchBucket:
						log.ErrorF("list error:%s, %v", s3.ErrCodeNoSuchBucket, aErr.Error())
					default:
						log.Error(aErr.Error())
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
				req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(info.AwsBucketInfo.Bucket),
					Key:    obj.Key,
				})
				if downloadUrl, e := req.Presign(5 * 3600 * time.Second); e == nil {
					fetchInfoChan <- object.FetchApiInfo{
						Bucket:  info.QiniuBucket,
						Key:     *obj.Key,
						FromUrl: downloadUrl,
					}
				} else {
					log.ErrorF("fetch([%s:%s]) create download url error: %v", info.AwsBucketInfo.Bucket, *obj.Key, e)
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
	waiter.Add(info.BatchInfo.WorkerCount)
	for i := 0; i < info.BatchInfo.WorkerCount; i++ {
		go func() {
			for fetchInfo := range fetchInfoChan {
				if _, e := object.Fetch(fetchInfo); e != nil {
					log.ErrorF("fetch %s => [%s:%s] failed, error:%v", fetchInfo.FromUrl, fetchInfo.Bucket, fetchInfo.Key, e)
					resultExport.Fail().ExportF("%s\t%s\t%v", fetchInfo.FromUrl, fetchInfo.Key, e)
				} else {
					log.AlertF("fetch %s => [%s:%s] success", fetchInfo.FromUrl, fetchInfo.Bucket, fetchInfo.Key)
					resultExport.Success().ExportF("%s\t%s", fetchInfo.FromUrl, fetchInfo.Key)
				}
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
}
