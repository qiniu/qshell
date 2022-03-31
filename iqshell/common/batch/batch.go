package batch

import "github.com/qiniu/qshell/v2/iqshell/common/flow"

type Info struct {
	flow.Info
	ExportInfo
	WorkCreatorInfo

	OverwriteExportFilePath string // 覆盖输出
}

func NewFlow() (*flow.Flow, error) {
	return nil, nil
}