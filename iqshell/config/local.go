package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/spf13/viper"
)

type ConfigType int8

const (
	ConfigTypeDefault ConfigType = 1
	ConfigTypeUser    ConfigType = 2
	ConfigTypeGlobal  ConfigType = 3
)

// local key
var (
	// qshell本地数据库文件目录
	localKeyAccountDBPath = []string{"path.accdb", "path.acc_db_path"}

	// qshell 账户文件目录
	localKeyAccountFilePath = []string{"path.acc", "path.acc_path"}

	// 上传HOST
	localKeyHostUp  = []string{"hosts.up", "hosts.up_host"}
	localKeyHostUps = []string{"hosts.ups", "hosts.up_hosts"}

	// RS HOST
	localKeyHostRs   = []string{"hosts.rs", "hosts.rs_host"}
	localKeyHostRses = []string{"hosts.rses", "hosts.rs_hosts"}

	// RSF HOST
	localKeyHostRsf  = []string{"hosts.rsf", "hosts.rsf_host"}
	localKeyHostRsfs = []string{"hosts.rsfs", "hosts.rsf_hosts"}

	// IO HOST
	localKeyHostIo  = []string{"hosts.io", "hosts.io_host"}
	localKeyHostIos = []string{"hosts.ios", "hosts.io_hosts"}

	// API HOST
	localKeyHostApi  = []string{"hosts.api", "hosts.api_host"}
	localKeyHostApis = []string{"hosts.apis", "hosts.api_hosts"}

	// UC HOST
	localKeyHostUc  = []string{"hosts.uc", "hosts.uc_host"}
	localKeyHostUcs = []string{"hosts.ucs", "hosts.uc_hosts"}

	// 账户密钥信息
	localKeyAccessKey = []string{"access_key"}
	localKeySecretKey = []string{"secret_key"}
)

var (
	userConfigViper   *viper.Viper
	globalConfigViper *viper.Viper
)

func GetAccountDBPath(configType ConfigType) string {
	return getStringValueFromLocal(configType, localKeyAccountDBPath)
}

func GetAccountFilePath(configType ConfigType) string {
	return getStringValueFromLocal(configType, localKeyAccountFilePath)
}

func GetUpHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostUp, localKeyHostUps)
}

func GetRsHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostRs, localKeyHostRses)
}

func GetRsfHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostRsf, localKeyHostRsfs)
}

func GetIoHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostIo, localKeyHostIos)
}

func GetUcHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostUc, localKeyHostUcs)
}

func GetApiHosts(configType ConfigType) []string {
	return getHostsFromLocal(configType, localKeyHostApi, localKeyHostApis)
}

func GetCredentials(configType ConfigType) auth.Credentials {
	return auth.Credentials{
		AccessKey: getAccessKey(configType),
		SecretKey: []byte(getSecretKey(configType)),
	}
}

func getAccessKey(configType ConfigType) string {
	return getStringValueFromLocal(configType, localKeyAccessKey)
}

func getSecretKey(configType ConfigType) string {
	return getStringValueFromLocal(configType, localKeySecretKey)
}

// ----- 业务封装
func getHostsFromLocal(configType ConfigType, hostKey []string, hostsKey []string) []string {
	var hosts []string
	hosts = getStringArrayValueFromLocal(configType, hostsKey)
	if hosts == nil {
		host := getStringValueFromLocal(configType, hostKey)
		if len(host) > 0 {
			hosts = []string{host}
		}
	}
	return hosts
}

func getStringValueFromLocal(configType ConfigType, localKey []string) string {
	value := ""
	vipers := getVipersWithConfigType(configType)
	if vipers != nil {
		for _, v := range vipers {
			value = getStringValueFromLocalByViper(v, localKey)
			if len(value) > 0 {
				break
			}
		}
	}
	return value
}

func getStringValueFromLocalByViper(viper *viper.Viper, localKey []string) string {
	value := ""
	for _, key := range localKey {
		value = viper.GetString(key)
		if len(value) > 0 {
			break
		}
	}
	return value
}

func getStringArrayValueFromLocal(configType ConfigType, localKey []string) []string {
	var value []string
	vipers := getVipersWithConfigType(configType)
	if vipers != nil {
		for _, viper := range vipers {
			value = getStringArrayValueFromLocalByViper(viper, localKey)
			if len(value) > 0 {
				break
			}
		}
	}
	return value
}

func getStringArrayValueFromLocalByViper(viper *viper.Viper, localKey []string) []string {
	var value []string
	for _, key := range localKey {
		value = viper.GetStringSlice(key)
		if len(value) > 0 {
			break
		}
	}
	return value
}

func getVipersWithConfigType(configType ConfigType) []*viper.Viper {
	var ret []*viper.Viper = nil
	addGlobalViper := func () {
		if globalConfigViper != nil {
			ret = append(ret, globalConfigViper)
		}
	}
	addUserViper := func () {
		if userConfigViper != nil {
			ret = append(ret, userConfigViper)
		}
	}
	switch configType {
	case ConfigTypeUser:
		addUserViper()
	case ConfigTypeGlobal:
		addGlobalViper()
	default:
		// 此顺序涉及优先级，不可调换
		addUserViper()
		addGlobalViper()
	}
	return ret
}