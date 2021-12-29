package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultDeadline = 3600
)

type PrivateUrlInfo struct {
	PublicUrl string
	Deadline  string
}

func (p PrivateUrlInfo) getDeadlineOfInt() (int64, error) {
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

type BatchPrivateUrlInfo struct {
	BatchInfo BatchInfo
	Deadline  string
}

// BatchPrivateUrl 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchPrivateUrl(info BatchPrivateUrlInfo) {
	if !prepareToBatch(info.BatchInfo) {
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
		if len(items) < 1 {
			continue
		}

		url := items[0]
		if url == "" {
			continue
		}

		urlToSign := strings.TrimSpace(url)
		if urlToSign == "" {
			continue
		}

		PrivateUrl(PrivateUrlInfo{
			PublicUrl: urlToSign,
			Deadline:  info.Deadline,
		})
	}
}
