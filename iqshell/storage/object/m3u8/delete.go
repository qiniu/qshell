package m3u8

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"sync"
)

type DeleteApiInfo struct {
	Bucket string
	Key    string
}

func Delete(info DeleteApiInfo) (err error) {
	m3u8FileList, err := Slices(SliceListApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
	})

	if err != nil {
		return errors.New("Get m3u8 file list error:" + err.Error())
	}

	if len(m3u8FileList) == 0 {
		return errors.New("no m3u8 slices found")
	}

	operations := make([]rs.BatchOperation,0, len(m3u8FileList))
	for _, file := range m3u8FileList {
		operations = append(operations, rs.DeleteApiInfo{
			Bucket:    file.Bucket,
			Key:       file.Key,
			AfterDays: 0,
		})
	}
	results, err := rs.Batch(operations)
	for result := range results {
		//TODO: 输出位置须再处理
		if result.Code != 200 || len(result.Error) > 0 {
			log.ErrorF("result error:%s", result.Error)
		}
	}

	return
}
