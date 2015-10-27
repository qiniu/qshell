package rs

import (
	"os"
	"testing"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
)

var (
	key        = "aa"
	newkey1    = "bbbb"
	newkey2    = "cccc"
	bucketName string
	domain     string
	client     Client
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	bucketName = os.Getenv("QINIU_TEST_BUCKET")
	domain = os.Getenv("QINIU_TEST_DOMAIN")
	if bucketName == "" || domain == "" {
		panic("require test env")
	}
	client = New(nil)

	// 删除 可能存在的 key
	client.Delete(nil, bucketName, key)
	client.Delete(nil, bucketName, newkey1)
	client.Delete(nil, bucketName, newkey2)
}

func upFile(localFile, bucketName, key string) error {

	policy := PutPolicy{
		Scope: bucketName + ":" + key,
	}
	return io.PutFile(nil, nil, policy.Token(nil), key, localFile, nil)
}

func TestEntry(t *testing.T) {

	//上传一个文件用用于测试
	err := upFile("rs_api.go", bucketName, key)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Delete(nil, bucketName, key)

	einfo, err := client.Stat(nil, bucketName, key)
	if err != nil {
		t.Fatal(err)
	}

	mime := "text/plain"
	err = client.ChangeMime(nil, bucketName, key, mime)
	if err != nil {
		t.Fatal(err)
	}

	einfo, err = client.Stat(nil, bucketName, key)
	if err != nil {
		t.Fatal(err)
	}
	if einfo.MimeType != mime {
		t.Fatal("mime type did not change")
	}

	err = client.Copy(nil, bucketName, key, bucketName, newkey1)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo, err := client.Stat(nil, bucketName, newkey1)
	if err != nil {
		t.Fatal(err)
	}
	if einfo.Hash != enewinfo.Hash {
		t.Fatal("invalid entryinfo:", einfo, enewinfo)
	}

	err = client.Move(nil, bucketName, newkey1, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo2, err := client.Stat(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	if enewinfo.Hash != enewinfo2.Hash {
		t.Fatal("invalid entryinfo:", enewinfo, enewinfo2)
	}

	err = client.Delete(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}

}
