package fop

import (
	"github.com/qiniu/rpc"
)

type ExifValType struct {
	Val  string `json:"val"`
	Type int    `json:"type"`
}

type ExifRet map[string]ExifValType

type Exif struct{}

func (this Exif) MakeRequest(url string) string {
	return url + "?exif"
}

func (this Exif) Call(l rpc.Logger, url string) (ret ExifRet, err error) {
	err = rpc.DefaultClient.Call(l, &ret, this.MakeRequest(url))
	return
}
