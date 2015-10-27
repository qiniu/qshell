package qshell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"qiniu/api.v6/conf"
	"qiniu/log"
)

type Account struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Zone      string `json:"zone,omitempty"`
}

func (this *Account) ToJson() (jsonStr string) {
	jsonStr = ""
	jsonData, err := json.Marshal(this)
	if err != nil {
		log.Error("Marshal account data failed.")
		return
	}
	jsonStr = string(jsonData)
	return jsonStr
}

func (this *Account) String() string {
	return fmt.Sprintf("AccessKey: %s SecretKey: %s Zone: %s", this.AccessKey, this.SecretKey, this.Zone)
}

func (this *Account) Set(accessKey string, secretKey string, zone string) {
	currentUser, err := user.Current()
	if err != nil {
		log.Error("Get current user failed.", err)
		return
	}
	qAccountFolder := filepath.Join(currentUser.HomeDir, ".qshell")
	if _, err := os.Stat(qAccountFolder); err != nil {
		if merr := os.MkdirAll(qAccountFolder, 0775); merr != nil {
			log.Error(fmt.Sprintf("Mkdir `%s' failed.", qAccountFolder), merr)
			return
		}
	}
	qAccountFile := filepath.Join(qAccountFolder, "account.json")
	log.Debug(fmt.Sprintf("Account file is `%s'", qAccountFile))
	fp, openErr := os.Create(qAccountFile)
	if openErr != nil {
		log.Error("Open account file failed.", openErr)
		return
	}
	defer fp.Close()

	this.AccessKey = accessKey
	this.SecretKey = secretKey
	account := Account{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Zone:      zone,
	}
	_, wErr := fp.WriteString(account.ToJson())
	if wErr != nil {
		log.Error("Write account info failed.", wErr)
		return
	}
}

func (this *Account) Get() {
	currentUser, err := user.Current()
	if err != nil {
		log.Error("Get current user failed.", err)
		return
	}
	qAccountFile := filepath.Join(currentUser.HomeDir, ".qshell", "account.json")
	fp, openErr := os.Open(qAccountFile)
	if openErr != nil {
		log.Error("Open account file failed.", openErr)
		return
	}
	defer fp.Close()
	accountBytes, readErr := ioutil.ReadAll(fp)
	if readErr != nil {
		log.Error("Read account file error.", readErr)
		return
	}
	if umError := json.Unmarshal(accountBytes, this); umError != nil {
		log.Error("Parse account file error.", umError)
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
}
