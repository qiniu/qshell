package utils

import "fmt"

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

// 转化文件大小到人工可读的字符串，以相应的单位显示
func FormatFileSize(size int64) (result string) {
	if size > TB {
		result = fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	} else if size > GB {
		result = fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	} else if size > MB {
		result = fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else if size > KB {
		result = fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	} else {
		result = fmt.Sprintf("%d B", size)
	}
	return
}