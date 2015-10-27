package gist

import (
	"log"

	"github.com/qiniu/api/fop"
)

// @gist makeImageInfoUrl
func makeImageInfoUrl(imageUrl string) string {
	ii := fop.ImageInfo{}
	return ii.MakeRequest(imageUrl)
}

// @endgist

func getImageInfo(imageUrl string) {

	var err error
	var ii = fop.ImageInfo{}
	var infoRet fop.ImageInfoRet

	// @gist getImageInfo
	infoRet, err = ii.Call(nil, imageUrl)
	if err != nil {
		// 产生错误
		log.Println("fop getImageInfo failed:", err)
		return
	}
	log.Println(infoRet.Height, infoRet.Width, infoRet.ColorModel,
		infoRet.Format)
	// @endgist
}

// @gist makeExifUrl
func makeExifUrl(imageUrl string) string {
	e := fop.Exif{}
	return e.MakeRequest(imageUrl)
}

// @endgist

func getExif(imageUrl string) {

	var err error
	var ie = fop.Exif{}
	var exifRet fop.ExifRet

	// @gist getExif
	exifRet, err = ie.Call(nil, imageUrl)
	if err != nil {
		// 产生错误
		log.Println("fop getExif failed:", err)
		return
	}

	// 处理返回结果
	for _, item := range exifRet {
		log.Println(item.Type, item.Val)
	}
	// @endgist
}

// @gist makeViewUrl
func makeViewUrl(imageUrl string) string {
	var view = fop.ImageView{
	// Mode    int    缩略模式
	// Width   int    Width = 0 表示不限定宽度
	// Height  int    Height = 0 表示不限定高度
	// Quality int    质量, 1-100
	// Format  string 输出格式，如jpg, gif, png, tif等等
	}
	return view.MakeRequest(imageUrl)
}

// @endgist
