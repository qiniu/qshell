package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func getStringValueFromLocal(vipers []*viper.Viper, localKey []string) *data.String {
	var value *data.String
	getValueFromLocal(vipers, localKey, func(v interface{}) (stop bool) {
		if v == nil {
			return false
		}
		value = data.NewString(cast.ToString(v))
		return true
	})
	return value
}

func getStringArrayValueFromLocal(vipers []*viper.Viper, localKey []string) []string {
	var value []string
	getValueFromLocal(vipers, localKey, func(v interface{}) (stop bool) {
		if v == nil {
			return false
		}
		value = cast.ToStringSlice(v)
		return true
	})
	return value
}

func getBoolValueFromLocal(vipers []*viper.Viper, localKey []string) *data.Bool {
	var value *data.Bool
	getValueFromLocal(vipers, localKey, func(v interface{}) (stop bool) {
		if v == nil {
			return false
		}
		value = data.NewBool(cast.ToBool(v))
		return true
	})
	return value
}

func getValueFromLocal(vipers []*viper.Viper, localKey []string, f func(value interface{}) (stop bool)) {
	if f == nil {
		return
	}

	stop := false
	var value interface{}
	for _, key := range localKey {
		for _, vip := range vipers {
			if value = vip.Get(key); value != nil {
				if stop = f(value); stop {
					break
				}
			}
		}

		if stop {
			break
		}
	}
}
