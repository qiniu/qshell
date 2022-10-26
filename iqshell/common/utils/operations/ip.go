package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils/ip"
	"time"
)

type IpQueryInfo struct {
	Ips []string
}

func (info *IpQueryInfo) Check() *data.CodeError {
	if len(info.Ips) == 0 {
		return alert.CannotEmptyError("Ip", "")
	}
	return nil
}

func IpQuery(cfg *iqshell.Config, info IpQueryInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if len(info.Ips) == 0 {
		log.Error(data.NewEmptyError().AppendDesc(alert.CannotEmpty("ip", "")))
		return
	}

	parser := ip.DefaultParser()
	for i, ipString := range info.Ips {
		if i > 0 {
			log.Alert("")
		}
		if result, err := parser.Parse(ipString); err != nil {
			log.Error(err)
			data.SetCmdStatusError()
		} else {
			log.AlertF("%v", result)
		}
		<-time.After(time.Millisecond * 500)
	}
}
