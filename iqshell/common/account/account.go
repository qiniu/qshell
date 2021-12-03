package account

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

// Name - 用户自定义的账户名称
type Account struct {
	Name      string
	AccessKey string
	SecretKey string
}

// 获取qbox.Mac
func (acc *Account) Mac() (mac *qbox.Mac) {

	mac = qbox.NewMac(acc.AccessKey, acc.SecretKey)
	return
}

// 对SecretKey进行加密， 保存AccessKey, 加密后的SecretKey在本地数据库中
func (acc *Account) Encrypt() (s string, err error) {
	encryptedKey, eErr := EncryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	s = strings.Join([]string{acc.Name, acc.AccessKey, encryptedKey}, ":")
	return
}

// 对SecretKey加密， 形成最后的数据格式
func (acc *Account) Value() (v string, err error) {
	encryptedKey, eErr := EncryptSecretKey(acc.AccessKey, acc.SecretKey)
	if eErr != nil {
		err = eErr
		return
	}
	v = Encrypt(acc.AccessKey, encryptedKey, acc.Name)
	return
}

// 保存在account.json文件中的数据格式
func Encrypt(accessKey, encryptedKey, name string) string {
	return strings.Join([]string{name, accessKey, encryptedKey}, ":")
}

func splits(joinStr string) []string {
	return strings.Split(joinStr, ":")
}

// 对保存在account.json中的文件字符串进行揭秘操作, 返回Account
func Decrypt(joinStr string) (acc Account, err error) {
	ss := splits(joinStr)
	name, accessKey, encryptedKey := ss[0], ss[1], ss[2]
	if name == "" || accessKey == "" || encryptedKey == "" {
		err = fmt.Errorf("name, accessKey and encryptedKey should not be empty")
		return
	}
	secretKey, dErr := DecryptSecretKey(accessKey, encryptedKey)
	if dErr != nil {
		err = fmt.Errorf("DecryptSecretKey: %v", dErr)
		return
	}
	acc.Name = name
	acc.AccessKey = accessKey
	acc.SecretKey = secretKey
	return
}

func (acc *Account) String() string {
	return fmt.Sprintf("Name: %s\nAccessKey: %s\nSecretKey: %s", acc.Name, acc.AccessKey, acc.SecretKey)
}

// 对SecretKey加密, 返回加密后的字符串
func EncryptSecretKey(accessKey, secretKey string) (string, error) {
	aesKey := utils.Md5Hex(accessKey)
	encryptedSecretKeyBytes, encryptedErr := utils.AesEncrypt([]byte(secretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return "", encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)
	return encryptedSecretKey, nil
}

// 对加密的SecretKey进行解密， 返回SecretKey
func DecryptSecretKey(accessKey, encryptedKey string) (string, error) {
	aesKey := utils.Md5Hex(accessKey)
	encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(encryptedKey)
	if decodeErr != nil {
		return "", decodeErr
	}
	secretKeyBytes, decryptErr := utils.AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
	if decryptErr != nil {
		return "", decryptErr
	}
	secretKey := string(secretKeyBytes)
	return secretKey, nil
}
