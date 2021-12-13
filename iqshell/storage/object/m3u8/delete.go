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

	handlerGroup := &sync.WaitGroup{}
	handlerGroup.Add(3)

	batchInfoChan := make(chan rs.DeleteApiInfo)
	go func() {
		for _, file := range m3u8FileList {
			batchInfoChan <- rs.DeleteApiInfo{
				Bucket:    file.Bucket,
				Key:       file.Key,
				AfterDays: 0,
			}
		}

		handlerGroup.Done()
	}()

	resultChan, errChan := rs.BatchDelete(batchInfoChan)
	go func() {
		err = <-errChan
		err = errors.New("batch error:" + err.Error())
		handlerGroup.Done()
	}()

	go func() {
		for result := range resultChan {
			//TODO: 输出位置须再处理
			if result.Code != 200 || len(result.Error) > 0 {
				log.ErrorF("result error:%s", result.Error)
			}
		}
		handlerGroup.Done()
	}()

	handlerGroup.Wait()
	return
}
