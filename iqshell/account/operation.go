package account

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/qshell/v2/iqshell/config"
	"github.com/qiniu/qshell/v2/iqshell/data"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 保存账户信息到账户文件中， 并保存在本地数据库
func SaveAccount(acc Account, accountOver bool) (err error) {
	sErr := SetAccountToLocalJson(acc)
	if sErr != nil {
		err = sErr
		return
	}

	err = SaveToDB(acc, accountOver)

	return
}

// 保存账户信息到账户文件中
func SetAccountToLocalJson(acc Account) (err error) {
	accountFh, openErr := os.OpenFile(info.accountPath, os.O_CREATE|os.O_RDWR, 0600)
	if openErr != nil {
		err = fmt.Errorf("Open account file error: %s", openErr)
		return
	}
	defer accountFh.Close()

	oldAccountFh, openErr := os.OpenFile(info.oldAccountPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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
	jsonStr, mErr := acc.Value()
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
	ldb, lErr := leveldb.OpenFile(info.accountDBPath, nil)
	if lErr != nil {
		err = fmt.Errorf("open db: %v", err)
		os.Exit(data.STATUS_HALT)
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
	ldbValue, mError := acc.Value()
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
		err = fmt.Errorf("Open account file error, %s, please use `account` to set AccessKey and SecretKey first", openErr)
		return
	}
	defer accountFh.Close()

	accountBytes, readErr := ioutil.ReadAll(accountFh)
	if readErr != nil {
		err = fmt.Errorf("Read account file error, %s", readErr)
		return
	}
	acc, dErr := Decrypt(string(accountBytes))
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
	return getAccount(info.oldAccountPath)
}

// 返回Account
func GetAccount() (account Account, err error) {
	ak, sk := config.AccessKey(), config.SecretKey()
	if ak != "" && sk != "" {
		return Account{
			AccessKey: ak,
			SecretKey: sk,
		}, nil
	}
	return getAccount(info.accountPath)
}

// 获取Mac
func GetMac() (mac *qbox.Mac, err error) {
	account, err := GetAccount()
	if err != nil {
		return nil, err
	}
	return account.Mac(), nil
}

// 切换账户
func ChUser(userName string) (err error) {
	if userName != "" {
		db, oErr := leveldb.OpenFile(info.accountDBPath, nil)
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
		user, dErr := Decrypt(string(value))
		if dErr != nil {
			err = fmt.Errorf("Decrypt account bytes: %v", dErr)
			return
		}

		return SetAccountToLocalJson(user)
	} else {
		rErr := os.Rename(info.oldAccountPath, info.accountPath+".tmp")
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.accountPath, info.oldAccountPath)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(info.accountPath+".tmp", info.accountPath)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}
	}
	return
}

// 获取用户列表
func GetUsers() (ret []*Account, err error) {

	db, gErr := leveldb.OpenFile(info.accountDBPath, nil)
	if gErr != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	var (
		value string
	)
	for iter.Next() {
		value = string(iter.Value())
		acc, dErr := Decrypt(value)
		if dErr != nil {
			err = fmt.Errorf("Decrypt account bytes: %v", dErr)
			return
		}
		ret = append(ret, &acc)
	}
	return
}

// 列举本地数据库记录的用户列表
func ListUser(userLsName bool) (err error) {
	db, gErr := leveldb.OpenFile(info.accountDBPath, nil)
	if gErr != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	var (
		name  string
		value string
	)
	for iter.Next() {
		name = string(iter.Key())
		value = string(iter.Value())
		acc, dErr := Decrypt(value)
		if dErr != nil {
			err = fmt.Errorf("Decrypt account bytes: %v", dErr)
			return
		}
		if userLsName {
			fmt.Println(name)
		} else {
			fmt.Printf("Name: %s\n", name)
			fmt.Printf("AccessKey: %s\n", acc.AccessKey)
			fmt.Printf("SecretKey: %s\n", acc.SecretKey)
			fmt.Println("")
		}
	}
	iter.Release()
	return
}

// 清除本地账户数据库
func CleanUser() (err error) {
	err = os.RemoveAll(info.accountDBPath)
	if err != nil {
		return
	}

	err = os.RemoveAll(info.accountPath)
	if err != nil {
		return
	}

	err = os.RemoveAll(info.oldAccountPath)

	return
}

// 从本地数据库删除用户
func RmUser(userName string) (err error) {
	if len(info.accountDBPath) == 0 {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, err := leveldb.OpenFile(info.accountDBPath, nil)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	err = db.Delete([]byte(userName), nil)
	logs.Debug("Removing user: %d\n", userName)
	return
}

// 查找用户
func LookUp(userName string) (err error) {
	if len(info.accountDBPath) == 0 {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, err := leveldb.OpenFile(info.accountDBPath, nil)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	var (
		name  string
		value string
	)
	for iter.Next() {
		name = string(iter.Key())
		value = string(iter.Value())
		acc, dErr := Decrypt(value)
		if dErr != nil {
			err = fmt.Errorf("Decrypt account bytes: %v", dErr)
			return
		}
		if strings.Contains(name, userName) {
			fmt.Println(acc.String())
		}
	}
	iter.Release()
	return
}