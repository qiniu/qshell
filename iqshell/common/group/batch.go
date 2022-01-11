package group

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/scanner"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
)

// Info Batch 参数
type Info struct {
	work.Info

	ItemSeparate           string
	InputFile              string // batch 操作输入文件
	Force                  bool   // 无需验证即可 batch 操作，类似于二维码验证
	Overwrite              bool   // 强制执行，服务端参数
	FailExportFilePath     string // 错误输出
	SuccessExportFilePath  string // 成功输出
	OverrideExportFilePath string // 覆盖输出
}

type Handler interface {
	Scanner() scanner.Scanner
	Export() *export.FileExporter
}

func NewHandler(info Info) (Handler, error) {
	if err := prepareToBatch(info); err != nil {
		return nil, err
	}

	e, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.SuccessExportFilePath,
		FailExportFilePath:     info.FailExportFilePath,
		OverrideExportFilePath: info.OverrideExportFilePath,
	})
	if err != nil {
		return nil, errors.New("get export error:" + err.Error())
	}

	s, err := scanner.NewScanner(scanner.Info{
		StdInEnable: true,
		InputFile:   info.InputFile,
	})
	if err != nil {
		return nil, errors.New("get scanner error:" + err.Error())
	}

	return &handler{
		export:  e,
		scanner: s,
	}, nil
}

func prepareToBatch(info Info) error {
	log.DebugF("forceFlag: %v, overwriteFlag: %v, worker: %v, inputFile: %q, bsuccessFname: %q, bfailureFname: %q, sep: %q",
		info.Force, info.Overwrite, info.WorkCount, info.InputFile, info.SuccessExportFilePath, info.FailExportFilePath, info.ItemSeparate)

	if !info.Force {
		return nil
	}

	code := utils.CreateRandString(6)
	log.Warning(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", code))

	confirm := ""
	_, err := fmt.Scanln(&confirm)
	if err != nil {
		return errors.New("scan error:" + err.Error())
	}

	if code != confirm {
		return errors.New("Task quit!")
	}

	return nil
}

type handler struct {
	export  *export.FileExporter
	scanner scanner.Scanner
}

func (b *handler) Scanner() scanner.Scanner {
	return b.scanner
}

func (b *handler) Export() *export.FileExporter {
	return b.export
}
