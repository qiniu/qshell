package batch

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
)

func Some(operations []Operation) ([]*OperationResult, *data.CodeError) {
	h := &someBatchHandler{
		readIndex:  0,
		operations: operations,
		results:    make([]*OperationResult, 0, len(operations)),
		err:        nil,
	}

	works := make([]flow.Work, len(operations))
	for _, operation := range operations {
		works = append(works, operation)
	}
	NewHandler(Info{
		Info: flow.Info{
			WorkerCount:       1,
			StopWhenWorkError: true,
		},
		WorkList:                    works,
		MaxOperationCountPerRequest: 1000,
	}).OnResult(func(operation Operation, result *OperationResult) {
		h.results = append(h.results, result)
	}).OnError(func(err *data.CodeError) {
		h.err = err
	}).Start()

	return h.results, h.err
}

type someBatchHandler struct {
	readIndex  int
	operations []Operation
	results    []*OperationResult
	err        *data.CodeError
}
