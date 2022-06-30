package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"strconv"
)

type FuncCallInfo struct {
	FuncTemplate string
	ParamsJson   string
	RunTimes     string
	runTimesInt  int
}

func (info *FuncCallInfo) Check() *data.CodeError {
	if len(info.FuncTemplate) == 0 {
		return alert.CannotEmptyError("FuncTemplate", "")
	}

	if len(info.ParamsJson) == 0 {
		return alert.CannotEmptyError("ParamsJson", "")
	}

	if len(info.RunTimes) == 0 {
		info.runTimesInt = 1
	} else {
		if times, err := strconv.Atoi(info.RunTimes); err != nil {
			return data.NewEmptyError().AppendDescF("Invalid RunTimes:%v", err)
		} else if times < 0 {
			info.runTimesInt = 0
		} else {
			info.runTimesInt = times
		}
	}
	return nil
}

func FuncCall(cfg *iqshell.Config, info FuncCallInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	t, tErr := utils.NewTemplate(info.FuncTemplate)
	if tErr != nil {
		log.ErrorF("%v", tErr)
		return
	}

	if output, err := t.RunWithJsonString(info.ParamsJson); err != nil {
		log.ErrorF("error:%v", err)
	} else {
		log.Warning("output is insert [], and you should be careful with spaces etc.")
		log.InfoF("[%s]", output)
	}
}
