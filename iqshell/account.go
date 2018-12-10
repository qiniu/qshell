package iqshell

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Account struct {
	Name      string
	AccessKey string
	SecretKey string
}

func (acc *Account) Mac() (mac *qbox.Mac) {

	mac = qbox.NewMac(acc.AccessKey, acc.SecretKey)
	return
}

func (acc *Account) Encrypt() (s string, err error) {
	encryptedKey, eErr := EncryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	s = strings.Join([]string{acc.Name, acc.AccessKey, encryptedKey}, ":")
	return
}

func (acc *Account) Value() (v string, err error) {
	encryptedKey, eErr := EncryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	v = Encrypt(acc.AccessKey, encryptedKey, acc.Name)
	return
}

func Encrypt(accessKey, encryptedKey, name string) string {
	return strings.Join([]string{name, accessKey, encryptedKey}, ":")
}

func splits(joinStr string) []string {
	return strings.Split(joinStr, ":")
}

func Decrypt(joinStr string) (acc Account, err error) {
	ss := splits(joinStr)
	name, accessKey, encryptedKey := ss[0], ss[1], ss[2]
	if name == "" || accessKey == "" || encryptedKey == "" {
		err = fmt.Errorf("name, accessKey and encryptedKey should not be empty")
		return
	}
	secretKey, dErr := DecryptSecretKey(accessKey, encryptedKey)
	if dErr != nil {
		err = fmt.Errorf("DecryptSecretKey: %v", dErr)
		return
	}
	acc.Name = name
	acc.AccessKey = accessKey
	acc.SecretKey = secretKey
	return
}

func (acc *Account) String() string {
	return fmt.Sprintf("Name: %s\nAccessKey: %s\nSecretKey: %s", acc.Name, acc.AccessKey, acc.SecretKey)
}

func EncryptSecretKey(accessKey, secretKey string) (string, error) {
	aesKey := Md5Hex(accessKey)
	encryptedSecretKeyBytes, encryptedErr := AesEncrypt([]byte(secretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return "", encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)
	return encryptedSecretKey, nil
}

func DecryptSecretKey(accessKey, encryptedKey string) (string, error) {
	aesKey := Md5Hex(accessKey)
	encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(encryptedKey)
	if decodeErr != nil {
		return "", decodeErr
	}
	secretKeyBytes, decryptErr := AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
	if decryptErr != nil {
		return "", decryptErr
	}
	secretKey := string(secretKeyBytes)
	return secretKey, nil
}

func setdb(acc Account, accountOver bool) (err error) {
	accDbPath := AccDBPath()
	if accDbPath == "" {
		return fmt.Errorf("empty account db path")
	}
	ldb, lErr := leveldb.OpenFile(accDbPath, nil)
	if lErr != nil {
		err = fmt.Errorf("open db: %v", err)
		os.Exit(STATUS_HALT)
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

func SetAccount2(accessKey, secretKey, name, accPath, oldPath string, accountOver bool) (err error) {
	acc := Account{
		Name:      name,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	sErr := SetAccount(acc, accPath, oldPath)
	if sErr != nil {
		err = sErr
		return
	}

	err = setdb(acc, accountOver)

	return
}

func SetAccount(acc Account, accPath, oldPath string) (err error) {
	QShellRootPath := RootPath()
	if QShellRootPath == "" {
		return fmt.Errorf("empty root path\n")
	}
	if _, sErr := os.Stat(QShellRootPath); sErr != nil {
		if mErr := os.MkdirAll(QShellRootPath, 0755); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` error: %s", QShellRootPath, mErr)
			return
		}
	}

	accountFh, openErr := os.OpenFile(accPath, os.O_CREATE|os.O_RDWR, 0600)
	if openErr != nil {
		err = fmt.Errorf("Open account file error: %s", openErr)
		return
	}
	defer accountFh.Close()

	oldAccountFh, openErr := os.OpenFile(oldPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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

func GetOldAccount() (account Account, err error) {
	AccountFname := OldAccPath()
	if AccountFname == "" {
		err = fmt.Errorf("empty old account path\n")
		return
	}

	return getAccount(AccountFname)
}

func GetAccount() (account Account, err error) {
	ak, sk := AccessKey(), SecretKey()
	if ak != "" && sk != "" {
		return Account{
			AccessKey: ak,
			SecretKey: sk,
		}, nil
	}
	AccountFname := AccPath()
	if AccountFname == "" {
		err = fmt.Errorf("empty account path\n")
		return
	}

	return getAccount(AccountFname)
}

func GetMac() (mac *qbox.Mac, err error) {
	account, err := GetAccount()
	if err != nil {
		return nil, err
	}
	return account.Mac(), nil
}

func ChUser(userName string) (err error) {
	if userName != "" {

		AccountDBPath := AccDBPath()
		if AccountDBPath == "" {
			err = fmt.Errorf("empty account db path\n")
			return
		}
		db, oErr := leveldb.OpenFile(AccountDBPath, nil)
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

		pt := AccPath()
		if pt == "" {
			err = fmt.Errorf("empty account path")
			return
		}
		oldPath := OldAccPath()
		if oldPath == "" {
			err = fmt.Errorf("empty account path")
			return
		}
		return SetAccount(user, pt, oldPath)
	} else {
		oldPath := OldAccPath()
		if oldPath == "" {
			err = fmt.Errorf("empty account path")
			return
		}
		pt := AccPath()
		if pt == "" {
			err = fmt.Errorf("empty account path")
			return
		}
		rErr := os.Rename(oldPath, pt+".tmp")
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}

		rErr = os.Rename(pt, oldPath)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}
		rErr = os.Rename(pt+".tmp", pt)
		if rErr != nil {
			err = fmt.Errorf("rename file: %v", rErr)
			return
		}
	}
	return
}

func GetUsers() (ret []*Account, err error) {

	AccountDBPath := AccDBPath()
	if AccountDBPath == "" {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, gErr := leveldb.OpenFile(AccountDBPath, nil)
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

func ListUser(userLsName bool) (err error) {
	AccountDBPath := AccDBPath()
	if AccountDBPath == "" {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, gErr := leveldb.OpenFile(AccountDBPath, nil)
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

func CleanUser() (err error) {
	QShellRootPath := RootPath()
	if QShellRootPath == "" {
		return fmt.Errorf("empty root path\n")
	}
	err = os.RemoveAll(QShellRootPath)
	return
}

func RmUser(userName string) (err error) {
	AccountDBPath := AccDBPath()
	if AccountDBPath == "" {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, err := leveldb.OpenFile(AccountDBPath, nil)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	db.Delete([]byte(userName), nil)
	logs.Debug("Removing user: %d\n", userName)
	return
}

func LookUp(userName string) (err error) {
	AccountDBPath := AccDBPath()
	if AccountDBPath == "" {
		err = fmt.Errorf("empty account db path\n")
		return
	}
	db, err := leveldb.OpenFile(AccountDBPath, nil)
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
