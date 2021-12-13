package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

var addCmdEg = ` qshell user add <AK> <SK> <UserName>
 or
 qshell user add --ak <AK> --sk <SK> --name <UserName>`

type AddInfo struct {
	Name      string
	AccessKey string
	SecretKey string
	Over      bool
}

// 保存账户信息到账户文件中， 并保存在本地数据库
func Add(info AddInfo) {
	if len(info.Name) == 0 {
		log.Error(alert.CannotEmpty("user name", addCmdEg))
		os.Exit(data.STATUS_ERROR)
		return
	}
	if len(info.AccessKey) == 0 {
		log.Error(alert.CannotEmpty("user ak", addCmdEg))
		os.Exit(data.STATUS_ERROR)
		return
	}
	if len(info.SecretKey) == 0 {
		log.Error(alert.CannotEmpty("user sk", addCmdEg))
		os.Exit(data.STATUS_ERROR)
		return
	}

	acc := account.Account{
		Name:      info.Name,
		AccessKey: info.AccessKey,
		SecretKey: info.SecretKey,
	}

	if err := account.SetAccountToLocalJson(acc); err != nil {
		log.ErrorF("user add: set current error:%v", err)
		os.Exit(data.STATUS_ERROR)
		return
	}

	if err := account.SaveToDB(acc, info.Over); err != nil {
		log.ErrorF("user add: save user to db error:%v", err)
		os.Exit(data.STATUS_ERROR)
		return
	}
}
