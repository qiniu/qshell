package batch

func One(operation Operation) (OperationResult, error) {
	results, err := Some([]Operation{operation})
	if len(results) == 0 {
		return OperationResult{}, err
	}
	return results[0], err
}
