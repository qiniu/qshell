package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"time"
)

type StatusInfo rs.StatusApiInfo

func Status(info StatusInfo) {
	result, err := rs.Status(rs.StatusApiInfo(info))
	if err != nil {
		log.ErrorF("Stat error:%v", err)
		return
	}

	log.Alert(getStatusInfo(info, result))
}

//func BatchStatus(info rs.StatusApiInfo) {
//	result, err := rs.Status(info)
//	if err != nil {
//		log.Error(err)
//		return
//	}
//
//	log.Alert(getStatusInfo(info, result))
//}

func getStatusInfo(info StatusInfo, status rs.StatusApiResult) string {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", info.Bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", info.Key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Hash:", status.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", status.FSize, utils.FormatFileSize(status.FSize))

	putTime := time.Unix(0, status.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", status.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", status.MimeType)
	if status.Type == 0 {
		statInfo += fmt.Sprintf("%-20s%d -> 标准存储\r\n", "FileType:", status.Type)
	} else {
		statInfo += fmt.Sprintf("%-20s%d -> 低频存储\r\n", "FileType:", status.Type)
	}
	return statInfo
}
