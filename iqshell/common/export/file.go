package export

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type FileExporter struct {
	success   Exporter
	fail      Exporter
	skip      Exporter
	overwrite Exporter
	result    Exporter
}

func (b *FileExporter) Success() Exporter {
	return b.success
}

func (b *FileExporter) Fail() Exporter {
	return b.fail
}

func (b *FileExporter) Skip() Exporter {
	return b.skip
}

func (b *FileExporter) Overwrite() Exporter {
	return b.overwrite
}

func (b *FileExporter) Result() Exporter {
	return b.result
}

func (b *FileExporter) Close() *data.CodeError {
	errS := b.success.Close()
	errF := b.fail.Close()
	errO := b.overwrite.Close()
	if errS == nil && errF == nil && errO == nil {
		return nil
	}
	return data.NewEmptyError().AppendDesc("export close:").
		AppendDesc("success").AppendError(errS).
		AppendDesc("fail").AppendError(errF).
		AppendDesc("overwrite").AppendError(errO)
}

type FileExporterConfig struct {
	SuccessExportFilePath   string // 输入列表中的成功部分
	FailExportFilePath      string // 输入列表中的失败部分
	SkipExportFilePath      string // 输入列表中的跳过部分
	OverwriteExportFilePath string // 输入列表中的覆盖部分
	ResultExportFilePath    string // 结果输出
}

func NewFileExport(config FileExporterConfig) (export *FileExporter, err *data.CodeError) {
	export = &FileExporter{}
	export.success, err = New(config.SuccessExportFilePath)
	if err != nil {
		return
	}

	export.fail, err = New(config.FailExportFilePath)
	if err != nil {
		return
	}

	export.skip, err = New(config.SkipExportFilePath)
	if err != nil {
		return
	}

	export.overwrite, err = New(config.OverwriteExportFilePath)
	if err != nil {
		return
	}

	export.result, err = New(config.ResultExportFilePath)
	return
}

func EmptyFileExport() *FileExporter {
	export := &FileExporter{}
	export.success = empty()
	export.fail = empty()
	export.skip = empty()
	export.overwrite = empty()
	export.result = empty()
	return export
}
