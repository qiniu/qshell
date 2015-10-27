package gist

import (
	"io"
	"log"

	"github.com/qiniu/api/rsf"
	"github.com/qiniu/rpc"

	. "github.com/qiniu/api/conf"
)

func init() {
	// @gist init
	ACCESS_KEY = "<YOUR_APP_ACCESS_KEY>"
	SECRET_KEY = "<YOUR_APP_SECRET_KEY>"
	// @endgist
}

// @gist listPrefix
func listAll(l rpc.Logger, rs *rsf.Client, bucketName string, prefix string) {

	var entries []rsf.ListItem
	var marker = ""
	var err error
	var limit = 1000

	for err == nil {
		entries, marker, err = rs.ListPrefix(l, bucketName,
			prefix, marker, limit)
		for _, item := range entries {
			// 处理 item
			log.Print("item:", item)
		}
	}
	if err != io.EOF {
		// 非预期的错误
		log.Print("listAll failed:", err)
	}
}

// @endgist
