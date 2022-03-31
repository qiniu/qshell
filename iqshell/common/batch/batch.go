package batch

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
)

type Info struct {
	flow.Info
	export.FileExporterConfig

	// 工作源
	WorkList    []flow.Work // 工作源：列表
	InputFile   string      // 工作源：文件
	EnableStdin bool        // 工作源：stdin, 当 InputFile 不存在时使用 stdin

	ItemSeparate            string // 分隔符
	Force                   bool   // 无需验证即可 batch 操作，类似于验证码验证
	OverwriteExportFilePath string // 覆盖输出
}

func NewHandler(info Info, creator flow.WorkCreator, workerProvider flow.WorkerProvider) (*Handler, error) {
	if creator == nil {
		return nil, alert.CannotEmptyError("batch handler: WorkCreator", "")
	}

	if workerProvider == nil {
		return nil, alert.CannotEmptyError("batch handler: WorkerProvider", "")
	}

	handler := &Handler{}
	if e, err := export.NewFileExport(info.FileExporterConfig); err != nil {
		return nil, err
	} else {
		handler.e = e
	}

	workProvider, err := flow.NewWorkProviderOfFile(info.InputFile, info.EnableStdin, creator)
	if err != nil {
		return nil, err
	}

	handler.f = &flow.Flow{
		Info:           info.Info,
		WorkProvider:   workProvider,
		WorkerProvider: workerProvider,
	}
	return handler, nil
}

type Handler struct {
	f *flow.Flow
	e *export.FileExporter
}

func (h *Handler) Flow() *flow.Flow {
	return h.f
}

func (h *Handler) Exporter() *export.FileExporter {
	return h.e
}
