package group

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
)

type Info struct {
	flow.Info
	export.FileExporterConfig

	// 工作数据源
	WorkList    []flow.Work // 工作数据源：列表
	InputFile   string      // 工作数据源：文件
	EnableStdin bool        // 工作数据源：stdin, 当 InputFile 不存在时使用 stdin

	// 解析每行数据
	ItemSeparate         string                                                     // 每行元素按分隔符分割：分隔符
	WorkBuilderWithItems func(items []string) (work flow.Work, err *data.CodeError) // 根据work 元素创建 work
	EnableParseJson      bool                                                       // 每行元素为 json：开启 json 检测
	BlankWorkBuilder     func() flow.Work                                           // 每行元素为 json：创建引用类型的 work，用于 json unmarshal

	// 数据处理
	WorkerProvider flow.WorkerProvider

	Force                   bool   // 无需验证即可 batch 操作，类似于验证码验证
	OverwriteExportFilePath string // 覆盖输出
}

func (info *Info) Check() *data.CodeError {
	if err := info.Info.Check(); err != nil {
		return err
	}

	if len(info.ItemSeparate) == 0 {
		info.ItemSeparate = data.DefaultLineSeparate
	}
	return nil
}

func NewHandler(info Info) (*Handler, *data.CodeError) {
	if err := info.Check(); err != nil {
		return nil, err
	}

	if info.WorkBuilderWithItems == nil && (!info.EnableParseJson || info.BlankWorkBuilder == nil) {
		return nil, alert.CannotEmptyError("batch handler: WorkBuilderWithItems", " set WorkBuilderWithItems or enable parse json")
	}

	if info.WorkerProvider == nil {
		return nil, alert.CannotEmptyError("batch handler: WorkerProvider", "")
	}

	var workCreator flow.WorkCreator
	if info.WorkBuilderWithItems != nil {
		workCreator = flow.NewLineSeparateWorkCreator(info.ItemSeparate, 0, info.WorkBuilderWithItems)
	} else {
		workCreator = flow.NewJsonWorkCreator(info.BlankWorkBuilder)
	}

	handler := &Handler{}
	if e, err := export.NewFileExport(info.FileExporterConfig); err != nil {
		return nil, err
	} else {
		handler.e = e
	}

	var err *data.CodeError
	var workProvider flow.WorkProvider
	if info.WorkList != nil && len(info.WorkList) > 0 {
		workProvider, err = flow.NewArrayWorkProvider(info.WorkList)
	} else {
		workProvider, err = flow.NewWorkProviderOfFile(info.InputFile, info.EnableStdin, workCreator)
	}
	if err != nil {
		return nil, err
	}

	handler.f = &flow.Flow{
		Info:           info.Info,
		WorkProvider:   workProvider,
		WorkerProvider: info.WorkerProvider,
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
