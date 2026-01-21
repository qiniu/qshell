package account

import (
	"fmt"
	"strings"

	"github.com/qiniu/qshell/v2/iqshell/common/data"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
)

// Account - 用户自定义的账户名称
type Account struct {
	Name      string
	AccessKey string
	SecretKey string
}

// 获取qbox.Mac
func (acc *Account) mac() (mac *qbox.Mac) {
	return qbox.NewMac(acc.AccessKey, acc.SecretKey)
}

// 对SecretKey进行加密， 保存AccessKey, 加密后的SecretKey在本地数据库中
func (acc *Account) encrypt() (s string, err *data.CodeError) {
	encryptedKey, eErr := encryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	s = strings.Join([]string{acc.Name, acc.AccessKey, encryptedKey}, ":")
	return
}

// 对SecretKey加密， 形成最后的数据格式
func (acc *Account) value() (v string, err *data.CodeError) {
	encryptedKey, eErr := encryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	v = encrypt(acc.AccessKey, encryptedKey, acc.Name)
	return
}

func (acc *Account) String() string {
	return fmt.Sprintf("Name: %s\nAccessKey: %s\nSecretKey: %s", acc.Name, acc.AccessKey, acc.SecretKey)
}
