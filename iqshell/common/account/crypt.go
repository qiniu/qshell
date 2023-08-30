package account

import (
	"encoding/base64"
	"strings"

	"github.com/qiniu/qshell/v2/iqshell/common/data"

	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

func splits(joinStr string) []string {
	return strings.Split(joinStr, ":")
}

// 对保存在account.json中的文件字符串进行揭秘操作, 返回Account
func decrypt(joinStr string) (acc Account, err *data.CodeError) {
	ss := splits(joinStr)
	if len(ss) != 3 {
		err = data.NewEmptyError().AppendDescF("account json style format error")
		return
	}

	name, accessKey, encryptedKey := ss[0], ss[1], ss[2]
	if name == "" || accessKey == "" || encryptedKey == "" {
		err = data.NewEmptyError().AppendDescF("name, accessKey and encryptedKey should not be empty")
		return
	}
	secretKey, dErr := decryptSecretKey(accessKey, encryptedKey)
	if dErr != nil {
		err = data.NewEmptyError().AppendDescF("DecryptSecretKey: %v", dErr)
		return
	}
	acc.Name = name
	acc.AccessKey = accessKey
	acc.SecretKey = secretKey
	return
}

// 保存在account.json文件中的数据格式
func encrypt(accessKey, encryptedKey, name string) string {
	return strings.Join([]string{name, accessKey, encryptedKey}, ":")
}

// 对SecretKey加密, 返回加密后的字符串
func encryptSecretKey(accessKey, secretKey string) (string, *data.CodeError) {
	aesKey := utils.Md5Hex(accessKey)
	encryptedSecretKeyBytes, encryptedErr := utils.AesEncrypt([]byte(secretKey), []byte(aesKey[7:23]))
	if encryptedErr != nil {
		return "", encryptedErr
	}
	encryptedSecretKey := base64.URLEncoding.EncodeToString(encryptedSecretKeyBytes)
	return encryptedSecretKey, nil
}

// 对加密的SecretKey进行解密， 返回SecretKey
func decryptSecretKey(accessKey, encryptedKey string) (string, *data.CodeError) {
	aesKey := utils.Md5Hex(accessKey)
	encryptedSecretKeyBytes, decodeErr := base64.URLEncoding.DecodeString(encryptedKey)
	if decodeErr != nil {
		return "", data.ConvertError(decodeErr)
	}
	secretKeyBytes, decryptErr := utils.AesDecrypt([]byte(encryptedSecretKeyBytes), []byte(aesKey[7:23]))
	if decryptErr != nil {
		return "", decryptErr
	}
	secretKey := string(secretKeyBytes)
	return secretKey, nil
}
