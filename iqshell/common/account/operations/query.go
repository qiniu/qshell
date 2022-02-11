package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

type ListInfo struct {
	OnlyListName bool
}

func (info *ListInfo) Check() error {
	return nil
}

func List(info ListInfo) {
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

// 当前用户
func Current() {
	acc, err := account.GetAccount()
	if err != nil {
		log.ErrorF("user current error: %v", err)
		os.Exit(data.StatusError)
	}
	log.AlertF(acc.String())
}

// 查找某个用户
type LookUpInfo struct {
	Name string
}

func (info *LookUpInfo)Check() error {
	if len(info.Name) == 0 {
		return alert.CannotEmptyError("user name", "")
	}
	return nil
}

func LookUp(info LookUpInfo) {
	acc, err := account.LookUp(info.Name)
	if err != nil {
		log.ErrorF("user lookup error: %v", err)
		os.Exit(data.StatusError)
	}
	log.AlertF(acc.String())
	return
}
