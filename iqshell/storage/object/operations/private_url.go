package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
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

	url, err := object.PrivateUrl(object.PrivateUrlApiInfo{
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
	handler, err := NewBatchHandler(info.BatchInfo)
	if err != nil {
		log.Error(err)
		return
	}

	work.NewFlowHandler(info.BatchInfo.Info.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, false
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) < 1 {
			return nil, true
		}
		url := items[0]
		if url == "" {
			return nil, true
		}
		urlToSign := strings.TrimSpace(url)
		if urlToSign == "" {
			return nil, true
		}
		return PrivateUrlInfo{
			PublicUrl: url,
			Deadline:  info.Deadline,
		}, true
	}).DoWork(func(work work.Work) (work.Result, error) {
		in := work.(PrivateUrlInfo)
		deadline, err := in.getDeadlineOfInt()
		if err != nil {
			return nil, err
		}
		return object.PrivateUrl(object.PrivateUrlApiInfo{
			PublicUrl: in.PublicUrl,
			Deadline:  deadline,
		})
	}).OnWorkResult(func(work work.Work, result work.Result) {
		url := work.(string)
		log.Alert(url)
	}).OnWorkError(func(work work.Work, err error) {
		log.Error(err)
	}).Start()
}
