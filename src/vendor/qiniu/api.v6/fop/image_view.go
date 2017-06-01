package fop

import (
	"strconv"
)

type ImageView struct {
	Mode    int    // 缩略模式
	Width   int    // Width = 0 表示不限定宽度
	Height  int    // Height = 0 表示不限定高度
	Quality int    // 质量, 1-100
	Format  string // 输出格式，如jpg, gif, png, tif等等
}

func (this ImageView) MakeRequest(url string) string {

	url += "?imageView/" + strconv.Itoa(this.Mode)
	if this.Width > 0 {
		url += "/w/" + strconv.Itoa(this.Width)
	}
	if this.Height > 0 {
		url += "/h/" + strconv.Itoa(this.Height)
	}
	if this.Quality > 0 {
		url += "/q/" + strconv.Itoa(this.Quality)
	}
	if this.Format != "" {
		url += "/format/" + this.Format
	}
	return url
}
