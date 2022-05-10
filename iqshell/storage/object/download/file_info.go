package download

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"path/filepath"
	"strings"
)

type fileInfo struct {
	toFile       string // 保存路径
	toAbsFile    string // 保存的绝对路径
	fileEncoding string // 文件编码方式

	fileDir   string // 保存文件的路径，从 ToFile 解析
	tempFile  string // 临时保存的文件路径 ToFile + .tmp
	fromBytes int64  // 下载开始位置，检查本地 tempFile 文件，读取已下载文件长度
}

func createDownloadFiles(toFile, fileEncoding string) (*fileInfo, *data.CodeError) {
	f := &fileInfo{
		toFile:       toFile,
		fileEncoding: fileEncoding,
	}

	err := f.check()
	if err != nil {
		return f, err
	}

	err = f.prepare()
	return f, err
}

func (d *fileInfo) check() *data.CodeError {
	if len(d.toFile) == 0 {
		return data.NewEmptyError().AppendDesc("the filename saved after downloading is empty")
	}
	return nil
}

func (d *fileInfo) prepare() *data.CodeError {
	// 文件路径
	var err error
	d.toAbsFile, err = filepath.Abs(d.toFile)
	if err != nil {
		err = data.NewEmptyError().AppendDesc("get save file abs path error:" + err.Error())
		return data.ConvertError(err)
	}

	if strings.ToLower(d.fileEncoding) == "gbk" {
		d.toAbsFile, err = utf82GBK(d.toFile)
		if err != nil {
			err = data.NewEmptyError().AppendDesc("gbk file path:" + d.toFile + " error:" + err.Error())
			return data.ConvertError(err)
		}
	}

	d.fileDir = filepath.Dir(d.toAbsFile)
	d.tempFile = fmt.Sprintf("%s.tmp", d.toAbsFile)

	err = os.MkdirAll(d.fileDir, 0775)
	if err != nil {
		return data.NewEmptyError().AppendDesc("MkdirAll failed for " + d.fileDir + " error:" + err.Error())
	}

	tempFileStatus, err := os.Stat(d.tempFile)
	if err != nil && os.IsNotExist(err) {
		d.fromBytes = 0
		return nil
	}

	if tempFileStatus != nil && !tempFileStatus.IsDir() {
		d.fromBytes = tempFileStatus.Size()
	}

	return nil
}

func (d *fileInfo) clean() *data.CodeError {
	err := os.Remove(d.toAbsFile)
	if e := os.Remove(d.tempFile); err == nil {
		err = e
	}
	return data.ConvertError(err)
}
