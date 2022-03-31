package batch

type ExportInfo struct {
	FailExportFilePath      string // 错误输出
	SuccessExportFilePath   string // 成功输出
	SkipExportFilePath      string // 跳过输出
	OverwriteExportFilePath string // 覆盖输出
}
