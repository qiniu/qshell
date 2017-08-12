package qshell

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Account struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
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

func SetAccount(accessKey string, secretKey string) (err error) {
	storageDir := filepath.Join(QShellRootPath, ".qshell")
	if _, sErr := os.Stat(storageDir); sErr != nil {
		if mErr := os.MkdirAll(storageDir, 0755); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` error, %s", storageDir, mErr)
			return
		}
	}

	accountFname := filepath.Join(storageDir, "account.json")

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

func GetAccount() (account Account, err error) {
	storageDir := filepath.Join(QShellRootPath, ".qshell")
	accountFname := filepath.Join(storageDir, "account.json")
	accountFh, openErr := os.Open(accountFname)
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

	logs.Debug("Load account from %s", accountFname)
	return
}
