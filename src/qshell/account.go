package qshell

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"qiniu/log"
)

type Account struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func (this *Account) ToJson() (jsonStr string, err error) {
	jsonData, mErr := json.Marshal(this)
	if mErr != nil {
		err = fmt.Errorf("Marshal account data failed, %s", mErr)
		return
	}
	jsonStr = string(jsonData)
	return
}

func (this *Account) String() string {
	return fmt.Sprintf("AccessKey: %s\nSecretKey: %s", this.AccessKey, this.SecretKey)
}

func (this *Account) Set(accessKey string, secretKey string) (err error) {
	qAccountFolder := ".qshell"
	if _, sErr := os.Stat(qAccountFolder); sErr != nil {
		if mErr := os.MkdirAll(qAccountFolder, 0775); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` failed, %s", qAccountFolder, mErr.Error())
			return
		}
	}
	qAccountFile := filepath.Join(qAccountFolder, "account.json")

	fp, openErr := os.Create(qAccountFile)
	if openErr != nil {
		err = fmt.Errorf("Open account file failed, %s", openErr.Error())
		return
	}
	defer fp.Close()

	aesKey := Md5Hex(accessKey)
	encryptedSecretKeyBytes, encryptedErr := AesEncrypt([]byte(secretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)

	this.AccessKey = accessKey
	this.SecretKey = encryptedSecretKey

	jsonStr, mErr := this.ToJson()
	if mErr != nil {
		err = mErr
		return
	}

	_, wErr := fp.WriteString(jsonStr)
	if wErr != nil {
		err = fmt.Errorf("Write account info failed, %s", wErr.Error())
		return
	}

	return
}

func (this *Account) Get() (err error) {
	qAccountFile := filepath.Join(".qshell", "account.json")
	fp, openErr := os.Open(qAccountFile)
	if openErr != nil {
		err = fmt.Errorf("Open account file failed, %s, please use `account` to set AccessKey and SecretKey first", openErr.Error())
		return
	}
	defer fp.Close()
	accountBytes, readErr := ioutil.ReadAll(fp)
	if readErr != nil {
		err = fmt.Errorf("Read account file error, %s", readErr.Error())
		return
	}
	if umError := json.Unmarshal(accountBytes, this); umError != nil {
		err = fmt.Errorf("Parse account file error, %s", umError.Error())
		return
	}

	// backwards compatible with old version of qshell, which encrypt ak/sk based on existing ak/sk
	if len(this.SecretKey) == 40 {
		this.Set(this.AccessKey, this.SecretKey)
		this.Get()
		return
	}

	aesKey := Md5Hex(this.AccessKey)
	encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(this.SecretKey)
	if decodeErr != nil {
		return decodeErr
	}
	secretKeyBytes, decryptErr := AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
	if decryptErr != nil {
		return decryptErr
	}
	this.SecretKey = string(secretKeyBytes)

	pwd, _ := os.Getwd()
	accountPath := filepath.Join(pwd, qAccountFile)
	log.Debugf("Load account from %s", accountPath)
	return
}
