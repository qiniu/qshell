package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type UserInfo struct{}

func (info *UserInfo) Check() *data.CodeError {
	return nil
}

func User(cfg *iqshell.Config, info UserInfo) {
	iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	})
}

type ListInfo struct {
	OnlyListName bool
}

func (info *ListInfo) Check() *data.CodeError {
	return nil
}

func List(cfg *iqshell.Config, info ListInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	accounts, err := account.GetUsers()
	if err != nil {
		log.ErrorF("user list error:%v", err)
		data.SetCmdStatusError()
		return
	}

	for index, acc := range accounts {
		if info.OnlyListName {
			log.Alert(acc.Name)
		} else {
			if index > 0 {
				log.Alert(" ")
			}
			log.Alert(acc.String())
		}
	}
}

// Current 当前用户
func Current(cfg *iqshell.Config) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: nil,
	}); !shouldContinue {
		return
	}

	acc, err := account.GetAccount()
	if err != nil {
		log.ErrorF("user current error: %v", err)
		data.SetCmdStatusError()
		return
	}
	log.Alert(acc.String())
}

// LookUpInfo 查找某个用户
type LookUpInfo struct {
	Name string
}

func (info *LookUpInfo) Check() *data.CodeError {
	if len(info.Name) == 0 {
		return alert.CannotEmptyError("user name", "")
	}
	return nil
}

func LookUp(cfg *iqshell.Config, info LookUpInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	accounts, err := account.LookUp(info.Name)
	if err != nil {
		log.ErrorF("user lookup error: %v", err)
		data.SetCmdStatusError()
		return
	}
	for _, acc := range accounts {
		log.Alert(acc.String())
		log.Alert("")
	}

	return
}
