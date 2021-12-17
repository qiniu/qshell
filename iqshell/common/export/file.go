package export

import "fmt"

type FileExporter struct {
	success  Exporter
	fail     Exporter
	override Exporter
}

func (b *FileExporter) Success() Exporter {
	return b.success
}

func (b *FileExporter) Fail() Exporter {
	return b.fail
}

func (b *FileExporter) Override() Exporter {
	return b.override
}

func (b *FileExporter) Close() error {
	errS := b.success.Close()
	errF := b.fail.Close()
	errO := b.override.Close()
	if errS == nil && errF == nil && errO == nil {
		return nil
	}
	return fmt.Errorf("export close: success error:%v fail error:%v override error:%v", errS, errF, errO)
}

type FileExporterConfig struct {
	SuccessExportFilePath  string
	FailExportFilePath     string
	OverrideExportFilePath string
}

func NewFileExport(config FileExporterConfig) (export *FileExporter, err error) {
	export = &FileExporter{}
	export.success, err = New(config.SuccessExportFilePath)
	if err != nil {
		return
	}

	export.fail, err = New(config.FailExportFilePath)
	if err != nil {
		return
	}

	export.override, err = New(config.OverrideExportFilePath)
	return
}
