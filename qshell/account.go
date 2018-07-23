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
	"strings"
)

type Account struct {
	Name      string
	AccessKey string
	SecretKey string
}

func (acc *Account) EncryptSecretKey() (string, error) {
	aesKey := Md5Hex(acc.AccessKey)
	encryptedSecretKeyBytes, encryptedErr := AesEncrypt([]byte(acc.SecretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return "", encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)
	return encryptedSecretKey, nil
}

func (acc *Account) DecryptSecretKey() (string, error) {
	aesKey := Md5Hex(acc.AccessKey)
	encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(acc.SecretKey)
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
	return fmt.Sprintf("Name: %s\nAccessKey: %s\nSecretKey: %s", acc.Name, acc.AccessKey, acc.SecretKey)
}

func setdb(acc Account) (err error) {
	db, err := sql.Open("sqlite3", AccountDBPath)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	sqlTable := `
	    create table if not exists user (
            uid integer primary key autoincrement,
            ak varchar(64),
            sk varchar(64),
	    name varchar(64)
	)
	`
	st, err := db.Prepare(sqlTable)
	if err != nil {
		return
	}
	_, err = st.Exec()
	if err != nil {
		return
	}
	err = insertAcc(acc, db)
	if err != nil {
		return
	}
	return db.Close()
}

func insertAcc(acc Account, db *sql.DB) (err error) {
	rows, err := db.Query("select * from user where ak=? and sk=?", acc.AccessKey, acc.SecretKey)
	if err != nil {
		return fmt.Errorf("select: %v\n", err)
	}
	var exists bool
	for rows.Next() {
		exists = true
	}
	rows.Close()
	if !exists {
		logs.Debug("insert user (%s, %s, %s)", acc.AccessKey, acc.SecretKey, acc.Name)
		st, err := db.Prepare("insert into user (ak, sk, name) values (?,?,?)")
		if err != nil {
			return err
		}
		_, err = st.Exec(acc.AccessKey, acc.SecretKey, acc.Name)
		if err != nil {
			return err
		}
	}
	return
}

func SetAccount2(accessKey, secretKey, name string) (err error) {
	if _, sErr := os.Stat(QShellRootPath); sErr != nil {
		if mErr := os.MkdirAll(QShellRootPath, 0755); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` error, %s", QShellRootPath, mErr)
			return
		}
	}
	accountFh, openErr := os.OpenFile(AccountFname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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
	account.Name = name

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
	err = setdb(account)

	return
}

func SetAccount(accessKey, secretKey string) (err error) {
	if _, sErr := os.Stat(QShellRootPath); sErr != nil {
		if mErr := os.MkdirAll(QShellRootPath, 0755); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` error, %s", QShellRootPath, mErr)
			return
		}
	}

	accountFh, openErr := os.OpenFile(AccountFname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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

func GetAccount() (account Account, err error) {

	accountFh, openErr := os.Open(AccountFname)
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

	if umError := json.Unmarshal(accountBytes, &account); umError != nil {
		err = fmt.Errorf("Parse account file error, %s", umError)
		return
	}

	// backwards compatible with old version of qshell, which encrypt ak/sk based on existing ak/sk
	if len(account.SecretKey) == 40 {
		setErr := SetAccount(account.AccessKey, account.SecretKey)
		if setErr != nil {
			return
		}
	} else {
		aesKey := Md5Hex(account.AccessKey)
		encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(account.SecretKey)
		if decodeErr != nil {
			err = decodeErr
			return
		}
		secretKeyBytes, decryptErr := AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
		if decryptErr != nil {
			err = decryptErr
			return
		}
		account.SecretKey = string(secretKeyBytes)
	}
	return
}

func ChUser(uid int) (err error) {
	db, err := sql.Open("sqlite3", AccountDBPath)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("select ak, sk, name from user where uid=?", uid)
	if err != nil {
		return fmt.Errorf("select: %v\n", err)
	}
	var user Account
	var exists bool
	for rows.Next() {
		exists = true
		rows.Scan(&user.AccessKey, &user.SecretKey, &user.Name)
		break
	}
	rows.Close()
	decrypted, err := user.DecryptSecretKey()
	if err != nil {
		return err
	}
	user.SecretKey = decrypted
	if !exists {
		fmt.Fprintf(os.Stderr, "account %s not exists\n", user.Name)
		os.Exit(1)
	}
	return SetAccount2(user.AccessKey, user.SecretKey, user.Name)
}

func ListUser() (err error) {
	db, err := sql.Open("sqlite3", AccountDBPath)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("select uid, name, ak, sk from user")
	if err != nil {
		return fmt.Errorf("select: %v\n", err)
	}
	var uid int
	var acc Account
	for rows.Next() {
		rows.Scan(&uid, &acc.Name, &acc.AccessKey, &acc.SecretKey)
		fmt.Printf("UID: %d\n", uid)
		fmt.Printf("Name: %s\n", acc.Name)
		fmt.Printf("AccessKey: %s\n", acc.AccessKey)
		secretKey, err := acc.DecryptSecretKey()
		if err != nil {
			return err
		}
		fmt.Printf("SecretKey: %s\n", secretKey)
		fmt.Println("")
	}
	return rows.Close()
}

func CleanUser() (err error) {
	err = os.RemoveAll(QShellRootPath)
	return
}

func RmUser(uid int) (err error) {
	db, err := sql.Open("sqlite3", AccountDBPath)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return
	}
	defer db.Close()
	st, err := db.Prepare("delete from user where uid=?")
	if err != nil {
		return err
	}
	logs.Debug("Removing user: %d\n", uid)
	st.Exec(uid)
	return st.Close()
}

func LookUp(userName string) error {
	db, err := sql.Open("sqlite3", AccountDBPath)
	if err != nil {
		err = fmt.Errorf("open db: %v", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("select uid, name, ak, sk from user")
	if err != nil {
		return fmt.Errorf("select: %v\n", err)
	}
	var uid int
	var acc Account
	for rows.Next() {
		rows.Scan(&uid, &acc.Name, &acc.AccessKey, &acc.SecretKey)
		if strings.Contains(acc.Name, userName) {
			fmt.Printf("UID: %d\n", uid)
			fmt.Printf("Name: %s\n", acc.Name)
			fmt.Printf("AccessKey: %s\n", acc.AccessKey)
			secretKey, err := acc.DecryptSecretKey()
			if err != nil {
				return err
			}
			fmt.Printf("SecretKey: %s\n", secretKey)
			fmt.Println("")
		}
	}
	return rows.Close()
}
