package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
	"time"
)

type PrivateUrlInfo struct {
	PublicUrl string
	Deadline  string
}

func (p PrivateUrlInfo)getDeadlineOfInt() (int64, error) {
	if len(p.Deadline) == 0 {
		return time.Now().Add(time.Second * 3600).Unix(), nil
	}

	if val, err := strconv.ParseInt(p.Deadline, 10, 64); err != nil {
		return 0, errors.New("invalid deadline")
	} else {
		return val, nil
	}
}

func PrivateUrl(info PrivateUrlInfo) {
	deadline, err := info.getDeadlineOfInt()
	if err != nil {
		log.Error(err)
		return
	}

	url, err := rs.PrivateUrl(rs.PrivateUrlApiInfo{
		PublicUrl: info.PublicUrl,
		Deadline:  deadline,
	})

	log.Alert(url)
}
