package docs

import _ "embed"

//go:embed batchcopy.md
var batchCopyDocument string

const BatchCopyType = "batchcopy"

func init() {
	addCmdDocumentInfo(BatchCopyType, batchCopyDocument)
}
