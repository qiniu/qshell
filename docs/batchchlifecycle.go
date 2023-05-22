package docs

import _ "embed"

//go:embed batchchlifecycle.md
var batchChangeLifecycleDocument string

const BatchChangeLifecycle = "batchchlifecycle"

func init() {
	addCmdDocumentInfo(BatchChangeLifecycle, batchChangeLifecycleDocument)
}
