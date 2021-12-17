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

//BatchFetch 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchFetch(info BatchFetchInfo) {
	if !prepareToBatch(info.BatchInfo) {
		return
	}

	resultExport, err := NewBatchResultExport(info.BatchInfo)
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	scanner, err := newBatchScanner(info.BatchInfo)
	if err != nil {
		log.ErrorF("get scanner error:%v", err)
		return
	}

	for {
		line, success := scanner.scanLine()
		if !success {
			break
		}

		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		key, fromUrl := "", ""
		if len(items) > 0 {
			fromUrl = items[0]
			if len(items) > 1 {
				key = items[1]
			} else if k, err := utils.KeyFromUrl(fromUrl); err == nil {
				key = k
			}
		}
		if key != "" && fromUrl != "" {
			_, err := rs.Fetch(rs.FetchApiInfo{
				Bucket:  info.Bucket,
				Key:     key,
				FromUrl: fromUrl,
			})

			if err != nil {
				resultExport.Fail().ExportF("%s\t%s\t%v", fromUrl, key, err)
				log.ErrorF("Fetch '%s' => %s:%s Failed, Error: %v", fromUrl, info.Bucket, key, err)
			} else {
				resultExport.Success().ExportF("%s\t%s", fromUrl, info.Bucket)
				log.AlertF("Fetch '%s' => %s:%s Success", fromUrl, info.Bucket, key)
			}
		}
	}
}
