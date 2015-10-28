package qshell

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
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
		err = errors.New(fmt.Sprintf("Marshal account data failed, %s", mErr))
		return
	}
	jsonStr = string(jsonData)
	return
}

func (this *Account) String() string {
	return fmt.Sprintf("AccessKey: %s SecretKey: %s Zone: %s", this.AccessKey, this.SecretKey, this.Zone)
}

func (this *Account) Set(accessKey string, secretKey string, zone string) (err error) {
	currentUser, uErr := user.Current()
	if uErr != nil {
		err = errors.New(fmt.Sprintf("Get current user failed, %s", uErr.Error()))
		return
	}
	qAccountFolder := filepath.Join(currentUser.HomeDir, ".qshell")
	if _, sErr := os.Stat(qAccountFolder); sErr != nil {
		if mErr := os.MkdirAll(qAccountFolder, 0775); mErr != nil {
			err = errors.New(fmt.Sprintf("Mkdir `%s' failed, %s", qAccountFolder, mErr.Error()))
			return
		}
	}
	qAccountFile := filepath.Join(qAccountFolder, "account.json")

	fp, openErr := os.Create(qAccountFile)
	if openErr != nil {
		err = errors.New(fmt.Sprintf("Open account file failed, %s", openErr.Error()))
		return
	}
	defer fp.Close()

	this.AccessKey = accessKey
	this.SecretKey = secretKey

	//default to nb
	if !IsValidZone(zone) {
		zone = ZoneNB
	}

	account := Account{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Zone:      zone,
	}

	jsonStr, mErr := account.ToJson()
	if mErr != nil {
		err = mErr
		return
	}

	_, wErr := fp.WriteString(jsonStr)
	if wErr != nil {
		err = errors.New(fmt.Sprintf("Write account info failed, %s", wErr.Error()))
		return
	}

	return
}

func (this *Account) Get() (err error) {
	currentUser, uErr := user.Current()
	if uErr != nil {
		err = errors.New(fmt.Sprintf("Get current user failed, %s", uErr.Error()))
		return
	}
	qAccountFile := filepath.Join(currentUser.HomeDir, ".qshell", "account.json")
	fp, openErr := os.Open(qAccountFile)
	if openErr != nil {
		err = errors.New(fmt.Sprintf("Open account file failed, %s", openErr.Error()))
		return
	}
	defer fp.Close()
	accountBytes, readErr := ioutil.ReadAll(fp)
	if readErr != nil {
		err = errors.New(fmt.Sprintf("Read account file error, %s", readErr.Error()))
		return
	}
	if umError := json.Unmarshal(accountBytes, this); umError != nil {
		err = errors.New(fmt.Sprintf("Parse account file error, %s", umError.Error()))
		return
	}

	if this.Zone == "" {
		this.Zone = ZoneNB
	}

	//set default hosts
	switch this.Zone {
	case ZoneAWS:
		SetZone(ZoneAWSConfig)
	case ZoneBC:
		SetZone(ZoneBCConfig)
	default:
		SetZone(ZoneNBConfig)
	}

	return
}
