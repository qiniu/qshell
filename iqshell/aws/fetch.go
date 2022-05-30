package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
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
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", cfg.CmdCfg.CmdId, info.AwsBucketInfo.Region, info.AwsBucketInfo.Bucket, info.QiniuBucket))
		return filepath.Join(cmdPath, jobId)
	}
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	exporter, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:   info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:      info.BatchInfo.FailExportFilePath,
		OverwriteExportFilePath: info.BatchInfo.OverwriteExportFilePath,
	})

	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	fetchInfoChan := make(chan flow.Work, info.BatchInfo.WorkerCount)
	// 生产者
	go func() {
		if e := listBucket(info.AwsBucketInfo, func(svc *s3.S3, obj *s3.Object) {
			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(info.AwsBucketInfo.Bucket),
				Key:    obj.Key,
			})
			if downloadUrl, e := req.Presign(5 * 3600 * time.Second); e == nil {
				fetchInfoChan <- &object.FetchApiInfo{
					Bucket:  info.QiniuBucket,
					Key:     *obj.Key,
					FromUrl: downloadUrl,
				}
				log.DebugF("get object:%s\t%d\t%s\t%s\n%s", *obj.Key, *obj.Size, *obj.ETag, *obj.LastModified, downloadUrl)
			} else {
				log.ErrorF("fetch([%s:%s]) create download url error: %v", info.AwsBucketInfo.Bucket, *obj.Key, e)
			}
		}); e != nil {
			log.Error(e)
		}
		close(fetchInfoChan)
	}()

	var overseer flow.Overseer
	if info.BatchInfo.EnableRecord {
		dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
		log.DebugF("aws batch fetch recorder:%s", dbPath)
		if overseer, err = flow.NewDBRecordOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: nil,
				},
				Result: &object.FetchResult{},
				Err:    nil,
			}
		}); err != nil {
			log.ErrorF("aws batch fetch create overseer error:%v", err)
			return
		}
	} else {
		log.Debug("aws batch fetch recorder:Not Enable")
	}

	metric := &batch.Metric{}
	metric.Start()
	flow.New(info.BatchInfo.Info).
		WorkProviderWithChan(fetchInfoChan).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*object.FetchApiInfo)
				return object.Fetch(*in)
			}), nil
		})).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		SetOverseer(overseer).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err == nil {
				return false, nil
			}

			if !info.BatchInfo.RecordRedoWhileError {
				return false, workRecord.Err
			}

			result, _ := workRecord.Result.(*object.FetchResult)
			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDesc("result is invalid")
			}
			return false, nil
		}).
		OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			operationResult, _ := result.(*object.FetchResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					log.DebugF("Skip line:%s because have done and success", work.Data)
				} else {
					metric.AddFailureCount(1)
					log.DebugF("Skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
				log.DebugF("Skip line:%s because:%v", work.Data, err)
			}

		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			metric.AddCurrentCount(1)
			metric.AddSuccessCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			in, _ := workInfo.Work.(*object.FetchApiInfo)
			exporter.Success().ExportF("%s\t%s", in.FromUrl, in.Bucket)
			log.InfoF("AWS Fetch Success, '%s' => [%s:%s]", in.FromUrl, in.Bucket, in.Key)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.AddFailureCount(1)
			metric.PrintProgress("AWS Batching:" + workInfo.Data)

			exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
			if in, ok := workInfo.Work.(*object.FetchApiInfo); ok {
				log.ErrorF("AWS Fetch Failed, '%s' => [%s:%s], Error: %v", in.FromUrl, in.Bucket, in.Key, err)
			} else {
				log.ErrorF("AWS Fetch Failed, %s, Error: %s", workInfo.Data, err)
			}
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save aws batch fetch result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save aws batch fetch result to path:%s", resultPath)
	}

	log.Info("------------- AWS Batch Result --------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("--------------------------------------------")
}
