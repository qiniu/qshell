package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
)

type DeleteInfo struct {
	Bucket    string
	Key       string
	AfterDays string
}

func (d DeleteInfo) getAfterDaysOfInt() (int, error) {
	if len(d.AfterDays) == 0 {
		return -1, nil
	}
	return strconv.Atoi(d.AfterDays)
}
func Delete(info DeleteInfo) {
	afterDays, err := info.getAfterDaysOfInt()
	if err != nil {
		log.ErrorF("delete after days invalid:%v", err)
		return
	}

	result, err := rs.Delete(rs.DeleteApiInfo{
		Bucket:    info.Bucket,
		Key:       info.Key,
		AfterDays: afterDays,
	})

	if err != nil {
		log.ErrorF("delete error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("delete error:%s", result.Error)
		return
	}
}
