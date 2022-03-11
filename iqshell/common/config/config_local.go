package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

	// UC HOST
	localKeyIsUseHttps = []string{"use_https"}

	// 账户密钥信息
	localKeyAccessKey = []string{"access_key"}
	localKeySecretKey = []string{"secret_key"}
)

var (
	userConfigViper   *viper.Viper
	globalConfigViper *viper.Viper
)

func GetAccountDBPath(configType ConfigType) string {
	return getStringValue(configType, localKeyAccountDBPath).Value()
}

func GetAccountFilePath(configType ConfigType) string {
	return getStringValue(configType, localKeyAccountFilePath).Value()
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
	return getStringValue(configType, localKeyAccessKey).Value()
}

func getSecretKey(configType ConfigType) string {
	return getStringValue(configType, localKeySecretKey).Value()
}

func getIsUseHttps(configType ConfigType) *data.Bool {
	return getBoolValue(configType, localKeyIsUseHttps)
}

func getStringValue(configType ConfigType, localKey []string) *data.String {
	return getStringValueFromLocal(getVipersWithConfigType(configType), localKey)
}

func getBoolValue(configType ConfigType, localKey []string) *data.Bool {
	return getBoolValueFromLocal(getVipersWithConfigType(configType), localKey)
}

func getHostsFromLocal(configType ConfigType, hostKey []string, hostsKey []string) []string {
	var hosts []string
	vipers := getVipersWithConfigType(configType)
	hosts = getStringArrayValueFromLocal(vipers, hostsKey)
	if hosts == nil {
		host := getStringValueFromLocal(vipers, hostKey)
		if data.NotEmpty(host) {
			hosts = []string{host.Value()}
		}
	}
	return hosts
}

func getVipersWithConfigType(configType ConfigType) []*viper.Viper {
	var ret []*viper.Viper = nil
	addGlobalViper := func() {
		if globalConfigViper != nil {
			ret = append(ret, globalConfigViper)
		}
	}
	addUserViper := func() {
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
