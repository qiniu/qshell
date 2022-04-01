package batch

import "github.com/qiniu/qshell/v2/iqshell/common/data"

func One(operation Operation) (*OperationResult, *data.CodeError) {
	results, err := Some([]Operation{operation})
	if len(results) == 0 {
		return nil, err
	}
	return results[0], err
}
