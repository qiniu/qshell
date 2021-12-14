package tools

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strconv"
	"time"
)

type TimestampInfo struct {
	Value string
}

// 转化unix时间戳为可读的字符串
func Timestamp2Date(info TimestampInfo) {
	ts, err := strconv.ParseInt(info.Value, 10, 64)
	if err != nil {
		log.ErrorF("Invalid timestamp Value:%s error:%s", info.Value, err)
		return
	}

	t := time.Unix(ts, 0)
	log.Alert(t.String())
}

// 转化毫秒时间戳到人工可读的字符串
func TimestampMilli2Date(info TimestampInfo) {
	tms, err := strconv.ParseInt(info.Value, 10, 64)
	if err != nil {
		log.ErrorF("Invalid mill timestamp Value:%s error:%s", info.Value, err)
		return
	}
	t := time.Unix(tms/1000, 0)
	log.Alert(t.String())
}

// 转化纳秒时间戳到人工可读的字符串, 百纳秒为单位，主要是对接七牛服务时间戳
func TimestampNano2Date(info TimestampInfo) {
	tns, err := strconv.ParseInt(info.Value, 10, 64)
	if err != nil {
		log.ErrorF("Invalid nano timestamp Value:%s error:%s", info.Value, err)
		return
	}
	t := time.Unix(0, tns*100)
	log.Alert(t.String())
}

// 转化时间字符串到unix时间戳
func Date2Timestamp(info TimestampInfo) {
	duration, err := strconv.ParseInt(info.Value, 10, 64)
	if err != nil {
		log.ErrorF("Invalid duration Value:%s error:%s", info.Value, err)
		return
	}

	t := time.Now()
	t = t.Add(time.Second * time.Duration(duration))
	log.Alert(t.String())
}