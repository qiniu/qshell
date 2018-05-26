package qshell

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Account struct {
	Uid        int
	Uname      string
	Primary    bool
	AccessKey  string
	SecretKey  string
	Updatetime int
}

func (acc *Account) ToJson() (jsonStr string, err error) {
	jsonData, mErr := json.Marshal(acc)
	if mErr != nil {
		err = fmt.Errorf("Marshal account data error, %s", mErr)
		return
	}
	jsonStr = string(jsonData)
	return
}

func (acc *Account) String() string {
	return fmt.Sprintf("AccessKey: %s\nSecretKey: %s", acc.AccessKey, acc.SecretKey)
}

func SetAccount(uid int, uname string, primary bool, accessKey string, secretKey string) (err error) {
	accountFname := QAccountFile
	if accountFname == "" {
		storageDir := filepath.Join(QShellRootPath, ".qshell")
		if _, sErr := os.Stat(storageDir); sErr != nil {
			if mErr := os.MkdirAll(storageDir, 0755); mErr != nil {
				err = fmt.Errorf("Mkdir `%s` error, %s", storageDir, mErr)
				return
			}
		}

		accountFname = filepath.Join(storageDir, "account.json")
	}

	accountFh, openErr := os.OpenFile(accountFname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if openErr != nil {
		err = fmt.Errorf("Open account file error, %s", openErr)
		return
	}
	defer accountFh.Close()

	//encrypt ak&sk
	aesKey := Md5Hex(accessKey)
	encryptedSecretKeyBytes, encryptedErr := AesEncrypt([]byte(secretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)

	//write to local dir
	var account Account
	account.AccessKey = accessKey
	account.SecretKey = encryptedSecretKey

	jsonStr, mErr := account.ToJson()
	if mErr != nil {
		err = mErr
		return
	}

	_, wErr := accountFh.WriteString(jsonStr)
	if wErr != nil {
		err = fmt.Errorf("Write account info error, %s", wErr)
		return
	}

	return
}

func GetAccount(uid int) (user Account, err error) {
	accountFname := QAccountFile
	if accountFname == "" {
		storageDir := filepath.Join(QShellRootPath, ".qshell")
		accountFname = filepath.Join(storageDir, "account.db")
	}
	if _, err := os.Stat(accountFname); os.IsNotExist(err) {
		err = fmt.Errorf("please use `account` to set AccessKey and SecretKey first")
		return
	}
	db, err = sql.Open("sqlite3", accountFname)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	sqlTable = `
	    create table if not exists user (
	    id integer
	    name varchar(64)
	    primary varchar(64)
	    updatetime datetime)
	`
	db.Exec(sqlTable)
	rows, err := db.Query("select id, name, primary, updatetime where id=?", uid)
	if err != nil {
		err = fmt.Errorf("select from table: %v", err)
		return
	}
	var user Account
	for rows.Next() {
		err = rows.Scan(&user.Uid, &user.Uname, &user.Primary, &user.Updatetime, &user.AccessKey, &user.SecretKey)
		if err != nil {
			err = fmt.Errorf("row scan: %s", err)
			return
		}
		statement, err := db.Prepare("update user set updatetime=? where id=?")
		if err != nil {
			err = fmt.Errorf("statement prepare: %s", err)
			return
		}
		statement.Exec(time.Now().Unix(), uid)
	}
	// backwards compatible with old version of qshell, which encrypt ak/sk based on existing ak/sk
	if len(user.SecretKey) == 40 {
		setErr := SetAccount(user.AccessKey, user.SecretKey)
		if setErr != nil {
			return
		}
	} else {
		aesKey := Md5Hex(user.AccessKey)
		encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(user.SecretKey)
		if decodeErr != nil {
			err = decodeErr
			return
		}
		secretKeyBytes, decryptErr := AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
		if decryptErr != nil {
			err = decryptErr
			return
		}
		user.SecretKey = string(secretKeyBytes)
	}

	logs.Info("Load account from %s", accountFname)
	return
}
