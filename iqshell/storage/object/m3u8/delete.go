package m3u8

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type DeleteApiInfo struct {
	Bucket string
	Key    string
}

func Delete(info DeleteApiInfo) ([]rs.OperationResult, error) {
	m3u8FileList, err := Slices(SliceListApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
	})

	if err != nil {
		return nil, errors.New("Get m3u8 file list error:" + err.Error())
	}

	if len(m3u8FileList) == 0 {
		return nil, errors.New("no m3u8 slices found")
	}

	operations := make([]rs.BatchOperation,0, len(m3u8FileList))
	for _, file := range m3u8FileList {
		operations = append(operations, rs.DeleteApiInfo{
			Bucket:    file.Bucket,
			Key:       file.Key,
			AfterDays: 0,
		})
	}

	return rs.Batch(operations)
}
