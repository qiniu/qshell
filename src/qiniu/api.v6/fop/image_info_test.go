package fop

import (
	"os"
	"testing"

	. "github.com/qiniu/api/conf"
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require test env")
	}
}

func TestImageInfo(t *testing.T) {
	info := ImageInfo{}
	ret, err := info.Call(nil, "http://cheneya.qiniudn.com/ffdfd_9")
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Format != "png" || ret.Width != 413 || ret.Height != 232 || ret.ColorModel != "nrgba" {
		t.Error("result not match")
	}
}
