package export

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type FileExporter struct {
	success   Exporter
	fail      Exporter
	skip      Exporter
	overwrite Exporter
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
	SuccessExportFilePath   string
	FailExportFilePath      string
	SkipExportFilePath      string
	OverwriteExportFilePath string
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
	return
}
