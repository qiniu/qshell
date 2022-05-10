package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

// 字符串的md5值
func Md5Hex(from string) string {
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(from))
	return hex.EncodeToString(md5Hasher.Sum(nil))
}

// 加密数据
func AesEncrypt(origData, key []byte) ([]byte, *data.CodeError) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, data.NewEmptyError().AppendError(err)
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 解密数据
func AesDecrypt(crypted, key []byte) ([]byte, *data.CodeError) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, data.NewEmptyError().AppendError(err)
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// 加密解密需要数据一定的格式， 如果愿数据不符合要求，需要加一些padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 加密解密需要数据一定的格式， 如果愿数据不符合要求，需要加一些padding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
