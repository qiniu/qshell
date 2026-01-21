package operations

import (
	"os"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type ZipInfo struct {
	ZipFilePath string
	UnzipPath   string
}

func (info *ZipInfo) Check() *data.CodeError {
	if len(info.ZipFilePath) == 0 {
		return alert.CannotEmptyError("QiniuZipFilePath", "")
	}
	return nil
}

// Unzip 解压使用mkzip压缩的文件
func Unzip(cfg *iqshell.Config, info ZipInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	var err error
	if len(info.UnzipPath) == 0 {
		info.UnzipPath, err = os.Getwd()
		if err != nil {
			data.SetCmdStatusError()
			log.Error("Get current work directory failed due to error", err)
			return
		}
	} else {
		if _, statErr := os.Stat(info.UnzipPath); statErr != nil {
			data.SetCmdStatusError()
			log.Error("Specified <UnzipToDir> is not a valid directory")
			return
		}
	}

	unzipErr := utils.Unzip(info.ZipFilePath, info.UnzipPath)
	if unzipErr != nil {
		data.SetCmdStatusError()
		log.Error("Unzip file failed due to error", unzipErr)
	}
}
