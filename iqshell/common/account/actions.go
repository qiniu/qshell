package account

import (
	"fmt"
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
func SetAccountToLocalFile(acc Account) (err error) {
	accountFh, openErr := os.OpenFile(info.AccountPath, os.O_CREATE|os.O_RDWR, 0600)
	if openErr != nil {
		err = fmt.Errorf("Open account file error: %s", openErr)
		return
	}
	defer accountFh.Close()

	oldAccountFh, openErr := os.OpenFile(info.OldAccountPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if openErr != nil {
		err = fmt.Errorf("Open account file error: %s", openErr)
		return
	}
	defer oldAccountFh.Close()

	_, cErr := io.Copy(oldAccountFh, accountFh)
	if cErr != nil {
		err = cErr
		return
	}
	jsonStr, mErr := acc.value()
	if mErr != nil {
		err = mErr
		return
	}
	_, sErr := accountFh.Seek(0, io.SeekStart)
	if sErr != nil {
		err = sErr
		return
	}
	tErr := accountFh.Truncate(0)
	if tErr != nil {
		err = tErr
		return
	}
	_, wErr := accountFh.WriteString(jsonStr)
	if wErr != nil {
		err = fmt.Errorf("Write account info error, %s", wErr)
		return
	}
	return
}

func SaveToDB(acc Account, accountOver bool) (err error) {
	ldb, lErr := leveldb.OpenFile(info.AccountDBPath, nil)
	if lErr != nil {
		err = fmt.Errorf("open db: %v", err)
		os.Exit(data.StatusHalt)
	}
	defer ldb.Close()

	if !accountOver {

		exists, hErr := ldb.Has([]byte(acc.Name), nil)
		if hErr != nil {
			err = hErr
			return
		}
		if exists {
			err = fmt.Errorf("Account Name: %s already exist in local db", acc.Name)
			return
		}
	}

	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}
	ldbValue, mError := acc.value()
	if mError != nil {
		err = fmt.Errorf("Account.Value: %v", mError)
		return
	}
	putErr := ldb.Put([]byte(acc.Name), []byte(ldbValue), &ldbWOpt)
	if putErr != nil {
		err = fmt.Errorf("leveldb Put: %v", putErr)
		return
	}
	return
}

func getAccount(pt string) (account Account, err error) {
	accountFh, openErr := os.Open(pt)
	if openErr != nil {
		err = fmt.Errorf("Open account file error, %s, please use `account` to set Id and SecretKey first", openErr)
		return
	}
	defer accountFh.Close()

	accountBytes, readErr := ioutil.ReadAll(accountFh)
	if readErr != nil {
		err = fmt.Errorf("Read account file error, %s", readErr)
		return
	}

	if len(accountBytes) == 0 {
		err = fmt.Errorf("Read account file error, account is empty")
		return
	}

	acc, dErr := decrypt(string(accountBytes))
	if dErr != nil {
		err = fmt.Errorf("Decrypt account bytes: %v", dErr)
		return
	}
	account = acc
	return
}

// qshell 会记录当前的user信息，当切换账户后， 老的账户信息会记录下来
// qshell user cu就可以切换到老的账户信息， 参考cd -回到先前的目录
func GetOldAccount() (account Account, err error) {
	return getAccount(info.OldAccountPath)
}

// 返回Account
func GetAccount() (account Account, err error) {
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
func GetMac() (mac *qbox.Mac, err error) {
	account, err := GetAccount()
	if err != nil {
		return nil, err
	}
	return account.mac(), nil
}

// 切换账户
func ChUser(userName string) (err error) {
	if userName != "" {
		db, oErr := leveldb.OpenFile(info.AccountDBPath, nil)
		if err != nil {
			err = fmt.Errorf("open db: %v", oErr)
			return
		}
		defer db.Close()

		value, gErr := db.Get([]byte(userName), nil)
		if gErr != nil {
			err = gErr
			return
		}
		user, dErr := decrypt(string(value))
		if dErr != nil {
			err = fmt.Errorf("Decrypt account bytes: %v", dErr)
			return
		}

		return SetAccountToLocalFile(user)
	} else {
		if _, err = GetOldAccount(); err != nil {
			return fmt.Errorf("get last account error:%v", err)
		}

		rErr := os.Rename(info.OldAccountPath, info.AccountPath+".tmp")
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.AccountPath, info.OldAccountPath)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.AccountPath+".tmp", info.AccountPath)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}
	}
	return
}

// 获取用户列表
func GetUsers() (ret []*Account, err error) {

	db, gErr := leveldb.OpenFile(info.AccountDBPath, nil)
	if gErr != nil {
		err = fmt.Errorf("open db: %v", err)
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
func CleanUser() (err error) {
	err = os.RemoveAll(info.AccountDBPath)
	if err != nil {
		return
	}

	err = os.RemoveAll(info.AccountPath)
	if err != nil {
		return
	}

	err = os.RemoveAll(info.OldAccountPath)

	return
}

// 从本地数据库删除用户
func RmUser(userName string) (err error) {
	if len(info.AccountDBPath) == 0 {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, err := leveldb.OpenFile(info.AccountDBPath, nil)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	err = db.Delete([]byte(userName), nil)
	log.DebugF("Removing user: %d\n", userName)
	return
}

// LookUp 查找用户
func LookUp(userName string) ([]Account, error) {
	if len(info.AccountDBPath) == 0 {
		return nil, fmt.Errorf("empty account db path\n")
	}

	db, err := leveldb.OpenFile(info.AccountDBPath, nil)
	if err != nil {
		return nil, fmt.Errorf("open db: %v", err)
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
