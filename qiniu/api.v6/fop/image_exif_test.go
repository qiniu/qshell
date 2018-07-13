package fop

import (
	"testing"
)

func TestImageExif(t *testing.T) {
	ie := Exif{}
	url := "http://cheneya.qiniudn.com/hello_jpg"
	_, err := ie.Call(nil, url)
	if err != nil {
		t.Error(err)
		return
	}
}
