package fop

import (
	"net/http"
	"os"
	"testing"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

var (
	key       = "gogopher.jpg"
	localFile = "gogopher.jpg"
	domain    string
	bucket    string
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	domain = os.Getenv("QINIU_TEST_DOMAIN")
	bucket = os.Getenv("QINIU_TEST_BUCKET")
	if ACCESS_KEY == "" || SECRET_KEY == "" || domain == "" || bucket == "" {
		panic("require test env")
	}
}

func makeUrl(key string) string {
	return "http://" + domain + "/" + key
}

func TestImageViewRequest(t *testing.T) {

	iv := ImageView{
		Mode:  1,
		Width: 250,
	}

	rawUrl := makeUrl(key)
	imageViewUrl := iv.MakeRequest(rawUrl)

	if imageViewUrl != rawUrl+"?imageView/1/w/250" {
		t.Error("result not match")
		return
	}

	iv.Mode = 2
	iv.Height = 250
	iv.Quality = 80
	iv.Format = "jpg"
	imageViewUrl = iv.MakeRequest(rawUrl)

	if imageViewUrl != rawUrl+"?imageView/2/w/250/h/250/q/80/format/jpg" {
		t.Error("result not match")
		return
	}
}

func TestPrivateImageView(t *testing.T) {

	//首先上传一个图片 用于测试
	policy := rs.PutPolicy{
		Scope: bucket + ":" + key,
	}
	err := io.PutFile(nil, nil, policy.Token(nil), key, localFile, nil)
	if err != nil {
		t.Errorf("TestPrivateImageView failed: %v", err)
		return
	}

	rawUrl := makeUrl(key)

	iv := ImageView{
		Mode:    2,
		Height:  250,
		Quality: 80,
	}
	imageViewUrl := iv.MakeRequest(rawUrl)
	p := rs.GetPolicy{}
	imageViewUrlWithToken := p.MakeRequest(imageViewUrl, nil)
	resp, err := http.DefaultClient.Get(imageViewUrlWithToken)
	if err != nil {
		t.Errorf("TestPrivateImageView failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if (resp.StatusCode / 100) != 2 {
		t.Errorf("TestPrivateImageView failed: resp.StatusCode = %v", resp.StatusCode)
		return
	}
}
