package io

import (
	"math/rand"
	"os"
	"testing"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

var (
	bucket   string
	testKey  = "resumableput_key"
	testFile = "resumable_api_test.go"
	mockerr  = false
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	bucket = os.Getenv("QINIU_TEST_BUCKET")
	if bucket == "" {
		panic("require QINIU_TEST_BUCKET")
	}
	rs.New(nil).Delete(nil, bucket, testKey)
}

func TestAll(t *testing.T) {

	policy := rs.PutPolicy{
		Scope: bucket,
	}
	token := policy.Token(nil)
	params := map[string]string{"x:1": "1"}
	extra := &PutExtra{
		ChunkSize: 128,
		MimeType:  "text/plain",
		Notify:    blockNotify,
		Params:    params,
	}

	testPut(t, token, nil)
	testPutWithoutKey(t, token, extra)
	testPutFile(t, token, extra)
	testPutFileWithoutKey(t, token, extra)
	testXVar(t, token, extra)
}

func testPut(t *testing.T, token string, extra *PutExtra) {

	var ret PutRet
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	err = Put(nil, &ret, token, testKey, f, fi.Size(), extra)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.New(nil).Delete(nil, bucket, ret.Key)
}

func testPutWithoutKey(t *testing.T, token string, extra *PutExtra) {

	var ret PutRet
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	err = PutWithoutKey(nil, &ret, token, f, fi.Size(), extra)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.New(nil).Delete(nil, bucket, ret.Key)
}

func testPutFile(t *testing.T, token string, extra *PutExtra) {

	var ret PutRet
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = PutFile(nil, &ret, token, testKey, testFile, extra)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.New(nil).Delete(nil, bucket, ret.Key)
}

func testPutFileWithoutKey(t *testing.T, token string, extra *PutExtra) {

	var ret PutRet
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = PutFileWithoutKey(nil, &ret, token, testFile, extra)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.New(nil).Delete(nil, bucket, ret.Key)
}

func testXVar(t *testing.T, token string, extra *PutExtra) {

	type Ret struct {
		PutRet
		X1 string `json:"x:1"`
	}
	var ret Ret
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	err = Put(nil, &ret, token, testKey, f, fi.Size(), extra)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.New(nil).Delete(nil, bucket, ret.Key)

	if ret.X1 != "1" {
		t.Fatal("test xVar failed:", ret.X1)
	}
}

//------------------------------------------------

func blockNotify(blkIdx int, blkSize int, ret *BlkputRet) {
	if rand.Int()%3 == 0 && mockerr == false {
		if ret.Ctx != "" {
			ret.Ctx = ""
			mockerr = true
		}
	}
}
