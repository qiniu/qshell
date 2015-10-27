package io

import (
	"bytes"
	"hash/crc32"
	"io"
	"os"
	"testing"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

var (
	bucket    string
	upString  = "hello qiniu world"
	policy    rs.PutPolicy
	localFile = "io_api.go"
	key1      = "test_put_1"
	key2      = "test_put_2"
	key3      = "test_put_3"
	extra     = []*PutExtra{
		&PutExtra{
			MimeType: "text/plain",
			CheckCrc: 0,
		},
		&PutExtra{
			MimeType: "text/plain",
			CheckCrc: 1,
		},
		&PutExtra{
			MimeType: "text/plain",
			CheckCrc: 2,
		},
		nil,
	}
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	bucket = os.Getenv("QINIU_TEST_BUCKET")
	if ACCESS_KEY == "" || SECRET_KEY == "" || bucket == "" {
		panic("require test env")
	}

	policy.Scope = bucket
}

//---------------------------------------

func crc32File(file string) uint32 {

	//it is so simple that do not check any err!!
	f, _ := os.Open(file)
	defer f.Close()
	info, _ := f.Stat()
	h := crc32.NewIEEE()
	buf := make([]byte, info.Size())
	io.ReadFull(f, buf)
	h.Write(buf)
	return h.Sum32()
}

func crc32String(s string) uint32 {

	h := crc32.NewIEEE()
	h.Write([]byte(s))
	return h.Sum32()
}

//---------------------------------------

func TestAll(t *testing.T) {

	testPut(t, key1)
	k1 := testPutWithoutKey(t)
	testPutFile(t, localFile, key2)
	k2 := testPutFileWithoutKey(t, localFile)

	testPut(t, key3)
	k3 := testPutWithoutKey2(t)

	//clear all keys
	rs.New(nil).Delete(nil, bucket, key1)
	rs.New(nil).Delete(nil, bucket, key2)
	rs.New(nil).Delete(nil, bucket, key3)
	rs.New(nil).Delete(nil, bucket, k1)
	rs.New(nil).Delete(nil, bucket, k2)
	rs.New(nil).Delete(nil, bucket, k3)
}

func testPut(t *testing.T, key string) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := Put(nil, ret, policy.Token(nil), key, buf, v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testPutWithoutKey(t *testing.T) string {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := PutWithoutKey(nil, ret, policy.Token(nil), buf, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	return ret.Key
}

func testPut2(t *testing.T, key string) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := Put2(nil, ret, policy.Token(nil), key, buf, int64(len(upString)), v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testPutWithoutKey2(t *testing.T) string {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := PutWithoutKey2(nil, ret, policy.Token(nil), buf, int64(len(upString)), v)
		if err != nil {
			t.Fatal(err)
		}
	}
	return ret.Key
}

func testPutFile(t *testing.T, localFile, key string) {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(localFile)
		}

		err := PutFile(nil, ret, policy.Token(nil), key, localFile, v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testPutFileWithoutKey(t *testing.T, localFile string) string {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(localFile)
		}

		err := PutFileWithoutKey(nil, ret, policy.Token(nil), localFile, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	return ret.Key
}
