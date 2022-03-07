package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

type UserInfo struct {
}

func (info *UserInfo) Check() error {
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

func (info *ListInfo) Check() error {
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
		log.ErrorF("user list error:", err)
		os.Exit(data.StatusError)
		return
	}

	for _, acc := range accounts {
		if info.OnlyListName {
			log.AlertF(acc.Name)
		} else {
			log.AlertF("Name: %s", acc.Name)
			log.AlertF("Id: %s", acc.AccessKey)
			log.AlertF("SecretKey: %s", acc.SecretKey)
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
		os.Exit(data.StatusError)
	}
	log.AlertF(acc.String())
}

// LookUpInfo 查找某个用户
type LookUpInfo struct {
	Name string
}

func (info *LookUpInfo) Check() error {
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

	acc, err := account.LookUp(info.Name)
	if err != nil {
		log.ErrorF("user lookup error: %v", err)
		os.Exit(data.StatusError)
	}
	log.AlertF(acc.String())
	return
}
