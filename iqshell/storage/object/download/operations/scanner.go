package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/scanner"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"strconv"
	"time"
)

type downloadScanner struct {
	exporter     *export.FileExporter
	lineScanner  scanner.Scanner
	itemSeparate string
	inputFile    string
	infoChan     chan *download.ApiInfo
	bucket       string
}

func newDownloadScanner(inputFile string, itemSeparate string, bucket string, exporter *export.FileExporter) (r *downloadScanner, err error) {
	r = &downloadScanner{
		exporter:     exporter,
		itemSeparate: itemSeparate,
		bucket:       bucket,
		inputFile:    inputFile,
		infoChan:     make(chan *download.ApiInfo, 100),
	}

	if len(inputFile) > 0 {
		r.lineScanner, err = scanner.NewScanner(scanner.Info{
			StdInEnable: true,
			InputFile:   inputFile,
		})

		if err != nil {
			err = errors.New("new download reader error:" + err.Error())
			return nil, err
		}

		// 配置输入文件则从输入文件中读取
		r.createFileScanOperation()
	} else {
		// 未配置输入文件则从 bucket 中读取
		r.createBucketScanOperation()
	}

	return
}

func (d *downloadScanner) createBucketScanOperation() {
	go func() {
		bucket.List(bucket.ListApiInfo{
			Bucket:    d.bucket,
			Prefix:    "",
			Marker:    "",
			Delimiter: "",
			StartTime: time.Time{},
			EndTime:   time.Time{},
			Suffixes:  nil,
			MaxRetry:  20,
		}, func(marker string, object bucket.ListObject) (bool, error) {
			d.infoChan <- &download.ApiInfo{
				Key:            object.Key,
				FileHash:       object.Hash,
				FileSize:       object.Fsize,
				FileModifyTime: object.PutTime,
			}
			return true, nil
		}, func(marker string, err error) {
		})
		close(d.infoChan)
	}()
}

func (d *downloadScanner) createFileScanOperation() {
	go func() {
		var keys []string
		for {
			if len(keys) == 100 {
				d.getDownloadObjectStatusAndAddToChan(keys)
				keys = nil
			}
			if keys == nil {
				keys = make([]string, 0, 100)
			}

			line, success := d.lineScanner.ScanLine()
			if !success {
				if len(keys) > 0 {
					d.getDownloadObjectStatusAndAddToChan(keys)
					keys = nil
				}
				break
			}

			items := utils.SplitString(line, d.itemSeparate)
			if len(items) < 1 || (len(items) > 0 && len(items[0]) == 0) {
				log.ErrorF("invalid line, line must contain key fileSize fileHash and fileModifyTime:%s", line)
				d.exporter.Fail().ExportF("%s: error:%s", line, "line must contain key fileSize fileHash and fileModifyTime")
				continue
			} else if len(items) < 4 {
				keys = append(keys, items[0])
				continue
			}

			fileKey := items[0]
			fileSize, err := strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				log.ErrorF("invalid line, get file size error:%s", line)
				d.exporter.Fail().ExportF("%s: get file size error:%s", line, err)
				continue
			}

			fileHash := items[2]
			fileModifyTime, err := strconv.ParseInt(items[3], 10, 64)
			if err != nil {
				log.ErrorF("invalid line, get file modify time error:%s", line)
				d.exporter.Fail().ExportF("%s: get file modify time error:%s", line, err)
				continue
			}

			d.infoChan <- &download.ApiInfo{
				Key:            fileKey,
				FileHash:       fileHash,
				FileSize:       fileSize,
				FileModifyTime: fileModifyTime,
			}
		}

		close(d.infoChan)
	}()
}

func (d *downloadScanner) getDownloadObjectStatusAndAddToChan(keys []string) {
	operations := make([]batch.Operation, 0, len(keys))
	for _, key := range keys {
		if len(key) > 0 {
			operations = append(operations, object.StatusApiInfo{
				Bucket: d.bucket,
				Key:    key,
			})
		}
	}
	results, err := batch.Some(operations)
	if err != nil {
		log.ErrorF("download batch status error:%v", err)
	}

	if len(results) == len(operations) {
		for i, result := range results {
			item := operations[i].(object.StatusApiInfo)
			if result.Code != 200 || result.Error != "" {
				d.exporter.Fail().ExportF("%s%s error:%v", item.Key, result.Error)
			} else {
				d.infoChan <- &download.ApiInfo{
					Key:            item.Key,
					FileHash:       result.Hash,
					FileSize:       result.FSize,
					FileModifyTime: result.PutTime,
				}
			}
		}
	}
}

func (d *downloadScanner) scan() (info *download.ApiInfo, hasMore bool) {
	for info = range d.infoChan {
		hasMore = true
		break
	}
	return
}

func (d *downloadScanner) getFileLineCount() int64 {
	if len(d.inputFile) == 0 {
		return 0
	}
	return utils.GetFileLineCount(d.inputFile)
}
