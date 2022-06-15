package account

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// 保存账户信息到账户文件中
func SetAccountToLocalFile(acc Account) (err *data.CodeError) {
	accountFh, openErr := os.OpenFile(info.AccountPath, os.O_CREATE|os.O_RDWR, 0600)
	if openErr != nil {
		err = data.NewEmptyError().AppendDescF("Open account file error: %s", openErr)
		return
	}
	defer accountFh.Close()

	oldAccountFh, openErr := os.OpenFile(info.OldAccountPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if openErr != nil {
		err = data.NewEmptyError().AppendDescF("Open account file error: %s", openErr)
		return
	}
	defer oldAccountFh.Close()

	_, cErr := io.Copy(oldAccountFh, accountFh)
	if cErr != nil {
		err = data.ConvertError(cErr)
		return
	}
	jsonStr, mErr := acc.value()
	if mErr != nil {
		err = mErr
		return
	}
	_, sErr := accountFh.Seek(0, io.SeekStart)
	if sErr != nil {
		err = data.ConvertError(sErr)
		return
	}
	tErr := accountFh.Truncate(0)
	if tErr != nil {
		err = data.ConvertError(tErr)
		return
	}
	_, wErr := accountFh.WriteString(jsonStr)
	if wErr != nil {
		err = data.NewEmptyError().AppendDescF("Write account info error, %s", wErr)
		return
	}
	return
}

func SaveToDB(acc Account, accountOver bool) (err *data.CodeError) {
	ldb, lErr := leveldb.OpenFile(info.AccountDBPath, nil)
	if lErr != nil {
		err = data.NewEmptyError().AppendDescF("open db: %v", err)
		os.Exit(data.StatusHalt)
	}
	defer ldb.Close()

	if !accountOver {

		exists, hErr := ldb.Has([]byte(acc.Name), nil)
		if hErr != nil {
			err = data.ConvertError(hErr)
			return
		}
		if exists {
			err = data.NewEmptyError().AppendDescF("Account Name: %s already exist in local db", acc.Name)
			return
		}
	}

	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}
	ldbValue, mError := acc.value()
	if mError != nil {
		err = data.NewEmptyError().AppendDescF("Account.Value: %v", mError)
		return
	}
	putErr := ldb.Put([]byte(acc.Name), []byte(ldbValue), &ldbWOpt)
	if putErr != nil {
		err = data.NewEmptyError().AppendDescF("leveldb Put: %v", putErr)
		return
	}
	return
}

func getAccount(pt string) (account Account, err *data.CodeError) {
	accountFh, openErr := os.Open(pt)
	if openErr != nil {
		err = data.NewEmptyError().AppendDescF("Open account file error, %s, please use `account` to set Id and SecretKey first", openErr)
		return
	}
	defer accountFh.Close()

	accountBytes, readErr := ioutil.ReadAll(accountFh)
	if readErr != nil {
		err = data.NewEmptyError().AppendDescF("Read account file error, %s", readErr)
		return
	}

	if len(accountBytes) == 0 {
		err = data.NewEmptyError().AppendDescF("Read account file error, account is empty")
		return
	}

	acc, dErr := decrypt(string(accountBytes))
	if dErr != nil {
		err = data.NewEmptyError().AppendDescF("Decrypt account bytes: %v", dErr)
		return
	}
	account = acc
	return
}

// qshell 会记录当前的user信息，当切换账户后， 老的账户信息会记录下来
// qshell user cu就可以切换到老的账户信息， 参考cd -回到先前的目录
func GetOldAccount() (account Account, err *data.CodeError) {
	return getAccount(info.OldAccountPath)
}

// 返回Account
func GetAccount() (account Account, err *data.CodeError) {
	credentials := config.GetCredentials(config.ConfigTypeDefault)
	if credentials.AccessKey != "" && credentials.SecretKey != nil {
		return Account{
			AccessKey: credentials.AccessKey,
			SecretKey: string(credentials.SecretKey),
		}, nil
	}
	return getAccount(info.AccountPath)
}

// 获取Mac
func GetMac() (mac *qbox.Mac, err *data.CodeError) {
	account, err := GetAccount()
	if err != nil {
		return nil, err
	}
	return account.mac(), nil
}

// 切换账户
func ChUser(userName string) (name string, err *data.CodeError) {
	if userName != "" {
		db, oErr := leveldb.OpenFile(info.AccountDBPath, nil)
		if oErr != nil {
			err = data.NewEmptyError().AppendDescF("open db: %v", oErr)
			return
		}
		defer db.Close()

		value, gErr := db.Get([]byte(userName), nil)
		if gErr != nil {
			err = data.NewEmptyError().AppendDescF("can't find user by name:%s , error:%v", userName, gErr)
			return
		}
		user, dErr := decrypt(string(value))
		if dErr != nil {
			err = data.NewEmptyError().AppendDescF("Decrypt account bytes: %v", dErr)
			return
		}

		return userName, SetAccountToLocalFile(user)
	} else {
		if acc, gErr := GetOldAccount(); gErr != nil {
			err = data.NewEmptyError().AppendDescF("get last account error:%v", gErr)
			return
		} else {
			name = acc.Name
		}

		rErr := os.Rename(info.OldAccountPath, info.AccountPath+".tmp")
		if rErr != nil {
			err = data.NewEmptyError().AppendDescF("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.AccountPath, info.OldAccountPath)
		if rErr != nil {
			err = data.NewEmptyError().AppendDescF("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.AccountPath+".tmp", info.AccountPath)
		if rErr != nil {
			err = data.NewEmptyError().AppendDescF("rename file: %v", rErr)
			return
		}
	}
	return
}

// 获取用户列表
func GetUsers() (ret []*Account, err *data.CodeError) {

	db, gErr := leveldb.OpenFile(info.AccountDBPath, nil)
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("open db: %v", err)
		return
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	var (
		name  string
		value string
	)
	for iter.Next() {
		name = string(iter.Key())
		value = string(iter.Value())
		acc, dErr := decrypt(value)
		if dErr != nil {
			log.WarningF("Decrypt account:%v error: %v", name, dErr)
			continue
		}
		ret = append(ret, &acc)
	}
	return
}

// 清除本地账户数据库
func CleanUser() *data.CodeError {
	err := os.RemoveAll(info.AccountDBPath)
	if err != nil {
		return data.ConvertError(err)
	}

	err = os.RemoveAll(info.AccountPath)
	if err != nil {
		return data.ConvertError(err)
	}

	err = os.RemoveAll(info.OldAccountPath)
	return data.ConvertError(err)
}

// 从本地数据库删除用户
func RmUser(userName string) *data.CodeError {
	if len(info.AccountDBPath) == 0 {
		return data.NewEmptyError().AppendDesc("empty account db path\n")
	}
	db, err := leveldb.OpenFile(info.AccountDBPath, nil)
	if err != nil {
		return data.NewEmptyError().AppendDescF("open db: %v", err)
	}
	defer db.Close()
	err = db.Delete([]byte(userName), nil)
	log.DebugF("Removing user: %s\n", userName)
	return nil
}

// LookUp 查找用户
func LookUp(userName string) ([]Account, *data.CodeError) {
	if len(info.AccountDBPath) == 0 {
		return nil, data.NewEmptyError().AppendDesc("empty account db path\n")
	}

	db, err := leveldb.OpenFile(info.AccountDBPath, nil)
	if err != nil {
		return nil, data.NewEmptyError().AppendDescF("open db: %v", err)
	}
	defer db.Close()

	accounts := make([]Account, 0, 1)
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	var name string
	var value string
	for iter.Next() {
		name = string(iter.Key())
		if strings.Contains(name, userName) {
			value = string(iter.Value())
			acc, DErr := decrypt(value)
			if DErr != nil {
				log.ErrorF("Decrypt account bytes: %v", err)
				continue
			}
			accounts = append(accounts, acc)
		}
	}
	return accounts, nil
}
