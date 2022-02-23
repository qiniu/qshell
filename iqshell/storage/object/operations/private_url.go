package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
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

func (info *PrivateUrlInfo) Check() error {
	if len(info.PublicUrl) == 0 {
		return alert.CannotEmptyError("PublicUrl", "")
	}
	return nil
}

func (p PrivateUrlInfo) getDeadlineOfInt() (int64, error) {
	if len(p.Deadline) == 0 {
		return time.Now().Add(time.Second * DefaultDeadline).Unix(), nil
	}

	if val, err := strconv.ParseInt(p.Deadline, 10, 64); err != nil {
		return 0, errors.New("invalid deadline")
	} else {
		return val, nil
	}
}

func PrivateUrl(cfg *iqshell.Config, info PrivateUrlInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	deadline, err := info.getDeadlineOfInt()
	if err != nil {
		log.Error(err)
		return
	}

	url, err := download.PublicUrlToPrivate(download.PublicUrlToPrivateApiInfo{
		PublicUrl: info.PublicUrl,
		Deadline:  deadline,
	})

	log.Alert(url)
}

type BatchPrivateUrlInfo struct {
	BatchInfo batch.Info
	Deadline  string
}

func (info *BatchPrivateUrlInfo) Check() error {
	return nil
}

// BatchPrivateUrl 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchPrivateUrl(cfg *iqshell.Config, info BatchPrivateUrlInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	handler, err := group.NewHandler(info.BatchInfo.Info)
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
		return download.PublicUrlToPrivate(download.PublicUrlToPrivateApiInfo{
			PublicUrl: in.PublicUrl,
			Deadline:  deadline,
		})
	}).OnWorkResult(func(work work.Work, result work.Result) {
		url := result.(string)
		log.Alert(url)
	}).OnWorkError(func(work work.Work, err error) {
		log.Error(err)
	}).Start()
}
