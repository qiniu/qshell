package rsf

import (
	"io"
	"os"
	"strconv"
	"testing"

	. "github.com/qiniu/api/conf"
	qio "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

var (
	bucketName string
	client     Client
	maxNum     = 1000
	keys       []string
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}

	bucketName = os.Getenv("QINIU_TEST_BUCKET")
	if bucketName == "" {
		panic("require test env")
	}
	client = New(nil)

}

func upFile(localFile, bucketName, key string) error {

	policy := rs.PutPolicy{
		Scope: bucketName + ":" + key,
	}
	return qio.PutFile(nil, nil, policy.Token(nil), key, localFile, nil)
}

func TestAll(t *testing.T) {

	//先上传文件到空间做初始化准备
	for i := 0; i < 10; i++ {
		key := "rsf_test_put_" + strconv.Itoa(i)
		err := upFile("rsf_api.go", bucketName, key)
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}
	defer func() {
		for _, k := range keys {
			rs.New(nil).Delete(nil, bucketName, k)
		}
	}()

	testList(t)
	testEof(t)
}
func testList(t *testing.T) {

	ret, marker, err := client.ListPrefix(nil, bucketName, "", "", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 5 && err != io.EOF {
		t.Fatal("TestList failed:", "expect len(ret) 5, but:", len(ret))
	}

	ret, _, err = client.ListPrefix(nil, bucketName, "", marker, 10000)
	if err != nil && err != io.EOF {
		t.Fatal("TestList failed:", "marker failed:", err)
	}
}

func testEof(t *testing.T) {

	_, _, err := client.ListPrefix(nil, bucketName, "", "", maxNum)

	if err != io.EOF {
		t.Fatal("TestEof failed:", "expect EOF but:", err)
	}

	_, _, err = client.ListPrefix(nil, bucketName, "", "", -1)

	if err != io.EOF {
		t.Fatal("TestEof failed:", "expect EOF but:", err)
	}
}
