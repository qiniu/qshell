package batch

import (
	"github.com/qiniu/qshell/v2/iqshell/common/work"
)

func Some(operations []Operation) ([]OperationResult, error) {
	handler := &someBatchHandler{
		readIndex:  0,
		operations: operations,
		results:    make([]OperationResult, 0, len(operations)),
		err:        nil,
	}

	NewFlow(Info{
		Info: work.Info{
			WorkCount:         1,
			StopWhenWorkError: true,
		},
		MaxOperationCountPerRequest: 1000,
	}).ReadOperation(func() (operation Operation, complete bool) {
		if handler.readIndex >= len(handler.operations) {
			return nil, true
		}
		return handler.operations[handler.readIndex], false
	}).OnResult(func(operation Operation, result OperationResult) {
		handler.results = append(handler.results, result)
	}).OnError(func(err error) {
		handler.err = err
	}).Start()

	return handler.results, handler.err
}

type someBatchHandler struct {
	readIndex  int
	operations []Operation
	results    []OperationResult
	err        error
}
