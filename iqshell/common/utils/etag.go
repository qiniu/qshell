package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	defaultChunkSize = 4 * 1024 * 1024
)

func GetEtag(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	return EtagV1(f)
}

func EtagV1(reader io.Reader) (string, error) {
	data, err := etagV1WithoutBase64Encoded(reader)
	if err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(data), nil
	}
}

func etagV1WithoutBase64Encoded(reader io.Reader) ([]byte, error) {
	var sha1s [][]byte
	for {
		dataBuf, err := ioutil.ReadAll(io.LimitReader(reader, defaultChunkSize))
		if err != nil {
			return nil, err
		} else if len(dataBuf) == 0 {
			break
		}
		sha1Result := sha1.Sum(dataBuf)
		sha1s = append(sha1s, sha1Result[:])
	}
	return hashSha1s(sha1s), nil
}

var ErrPartSizeMismatch = errors.New("part size mismatch with data from reader")

func EtagV2(reader io.Reader, parts []int64) (string, error) {
	if is4MbParts(parts) {
		return EtagV1(reader)
	}

	var sha1Buf []byte
	for _, partSize := range parts {
		dataBuf, err := ioutil.ReadAll(io.LimitReader(reader, partSize))
		if err != nil {
			return "", err
		} else if len(dataBuf) != int(partSize) {
			return "", ErrPartSizeMismatch
		}
		etagResult, err := etagV1WithoutBase64Encoded(bytes.NewReader(dataBuf))
		if err != nil {
			return "", err
		}
		sha1Buf = append(sha1Buf, etagResult[1:]...)
	}
	sha1Result := sha1.Sum(sha1Buf)
	sha1Buf = sha1Buf[:0]
	sha1Buf = append(sha1Buf, 0x9e)
	sha1Buf = append(sha1Buf, sha1Result[:]...)
	return base64.URLEncoding.EncodeToString(sha1Buf), nil
}

func IsSignByEtagV2(etag string) bool {
	etagData, err := base64.URLEncoding.DecodeString(etag)
	if len(etagData) < 1 || err != nil {
		return false
	} else {
		return etagData[0] == byte(0x9e)
	}
}

func hashSha1s(sha1s [][]byte) []byte {
	switch len(sha1s) {
	case 0:
		return []byte{
			0x16, 0xda, 0x39, 0xa3, 0xee, 0x5e, 0x6b,
			0x4b, 0x0d, 0x32, 0x55, 0xbf, 0xef, 0x95,
			0x60, 0x18, 0x90, 0xaf, 0xd8, 0x07, 0x09}
	case 1:
		var sha1Buf []byte
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf = append(sha1Buf, sha1s[0]...)
		return sha1Buf
	default:
		var (
			sha1Buf []byte
			buf     []byte
		)
		for _, sha1 := range sha1s {
			sha1Buf = append(sha1Buf, sha1...)
		}
		buf = append(buf, 0x96)
		sha1Result := sha1.Sum(sha1Buf)
		buf = append(buf, sha1Result[:]...)
		return buf
	}
}

func is4MbParts(parts []int64) bool {
	for _, partSize := range parts[:len(parts)-1] {
		if partSize != defaultChunkSize {
			return false
		}
	}
	if parts[len(parts)-1] > defaultChunkSize {
		return false
	}
	return true
}
