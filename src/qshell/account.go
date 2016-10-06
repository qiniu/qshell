package qshell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Account struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Zone      string `json:"zone,omitempty"`
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
	return fmt.Sprintf("AccessKey: %s SecretKey: %s Zone: %s", this.AccessKey, this.SecretKey, this.Zone)
}

func (this *Account) Set(accessKey string, secretKey string, zone string) (err error) {
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

	this.AccessKey = accessKey
	this.SecretKey = secretKey

	//default to nb
	if !IsValidZone(zone) {
		zone = ZoneNB
	}

	this.Zone = zone

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

	if this.Zone == "" {
		this.Zone = ZoneNB
	}

	//set default hosts
	SetZone(this.Zone)

	return
}
