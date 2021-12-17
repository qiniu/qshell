package operations

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

// BatchInfo Batch 参数
type BatchInfo struct {
	ItemSeparate           string
	InputFile              string
	Force                  bool // 无需验证即可 batch 操作，类似于二维码验证
	Overwrite              bool // 强制执行，服务端参数
	Worker                 int
	FailExportFilePath     string
	SuccessExportFilePath  string
	OverrideExportFilePath string
}

func prepareToBatch(info BatchInfo) bool {
	log.DebugF("forceFlag: %v, overwriteFlag: %v, worker: %v, inputFile: %q, bsuccessFname: %q, bfailureFname: %q, sep: %q",
		info.Force, info.Overwrite, info.Worker, info.InputFile, info.SuccessExportFilePath, info.FailExportFilePath, info.ItemSeparate)

	if !info.Force {
		return true
	}

	code := utils.CreateRandString(6)
	log.Warning(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", code))

	confirm := ""
	_, err := fmt.Scanln(&confirm)
	if err != nil {
		log.Error("scan error:" + err.Error())
		return false
	}

	if code != confirm {
		log.Error("Task quit!")
		return false
	}
	return true
}

// 输入
func newBatchScanner(info BatchInfo) (*batchScanner, error) {
	s := &batchScanner{}
	if len(info.InputFile) == 0 {
		s.scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(info.InputFile)
		if err != nil {
			return nil, errors.New("open src dest key map file error")
		}
		s.file = f
		s.scanner = bufio.NewScanner(f)
	}
	return s, nil
}

type batchScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (b *batchScanner) scanLine() (line string, success bool) {
	success = b.scanner.Scan()
	if success {
		line = b.scanner.Text()
	}
	log.DebugF("scan line:%s success:%t", line, success)
	return
}

func (b *batchScanner) close() error {
	if b.file == nil {
		return nil
	}
	return b.file.Close()
}
