package docs

import _ "embed"

//go:embed chlifecycle.md
var changeLifecycleDocument string

const ChangeLifecycle = "chlifecycle"

func init() {
	addCmdDocumentInfo(ChangeLifecycle, changeLifecycleDocument)
}
