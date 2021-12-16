package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
)

type FetchInfo rs.FetchApiInfo

func Fetch(info FetchInfo) {
	result, err := rs.Fetch(rs.FetchApiInfo(info))
	if err != nil {
		log.ErrorF("Fetch error: %v", err)
		os.Exit(data.STATUS_ERROR)
	} else {
		log.AlertF("Key:%s", result.Key)
		log.AlertF("Hash:%s", result.Hash)
		log.AlertF("Fsize: %d (%s)", result.Fsize, utils.FormatFileSize(result.Fsize))
		log.AlertF("Mime:%s", result.MimeType)
	}
}

type BatchFetchInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

// BatchFetch 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
//func BatchFetch(info BatchFetchInfo) {
//	if !prepareToBatch(info.BatchInfo) {
//		return
//	}
//
//	resultExport, err := NewBatchResultExport(info.BatchInfo)
//	if err != nil {
//		log.ErrorF("get export error:%v", err)
//		return
//	}
//
//	scanner, err := newBatchScanner(info.BatchInfo)
//	if err != nil {
//		log.ErrorF("get scanner error:%v", err)
//		return
//	}
//
//	rs.BatchWithHandler(&batchFetchHandler{
//		scanner:      scanner,
//		info:         &info,
//		resultExport: resultExport,
//	})
//
//}
//
//type batchFetchHandler struct {
//	scanner      *batchScanner
//	info         *BatchFetchInfo
//	resultExport *BatchResultExport
//}
//
//func (b batchFetchHandler) WorkCount() int {
//	return b.info.BatchInfo.Worker
//}
//
//func (b batchFetchHandler) ReadOperation() (rs.BatchOperation, bool) {
//	var info rs.BatchOperation = nil
//
//	line, success := b.scanner.scanLine()
//	if !success {
//		return nil, true
//	}
//
//	items := utils.SplitString(line, b.info.BatchInfo.ItemSeparate)
//	if len(items) > 0 {
//		key := ""
//		fromUrl := items[0]
//		if len(items) > 1 {
//			key = items[1]
//		} else if k, err := utils.KeyFromUrl(fromUrl); err == nil {
//			key = k
//		}
//
//		if key != "" && fromUrl != "" {
//			info = rs.FetchApiInfo{
//				Bucket:  b.info.Bucket,
//				Key:     key,
//				FromUrl: fromUrl,
//			}
//		}
//	}
//
//	return info, false
//}
//
//func (b batchFetchHandler) HandlerResult(operation rs.BatchOperation, result rs.OperationResult) {
//	apiInfo, ok := (operation).(rs.FetchApiInfo)
//	if !ok {
//		return
//	}
//
//	if result.Code != 200 || result.Error != "" {
//		b.resultExport.Fail().ExportF("%s\t%s\t%d\t%s", apiInfo.FromUrl, apiInfo.Key, result.Code, result.Error)
//		log.ErrorF("Fetch '%s' => %s:%s Failed, Code: %d, Error: %s", apiInfo.FromUrl, apiInfo.Bucket, apiInfo.Key, result.Code, result.Error)
//	} else {
//		b.resultExport.Success().ExportF("%s\t%s", apiInfo.FromUrl, apiInfo.Bucket)
//		log.AlertF("Fetch '%s' => %s:%s Success", apiInfo.FromUrl, apiInfo.Bucket, apiInfo.Key)
//	}
//}
//
//func (b batchFetchHandler) HandlerError(err error) {
//	log.ErrorF("batch Fetch error:%v:", err)
//}
