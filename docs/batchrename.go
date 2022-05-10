package docs

import _ "embed"

//go:embed batchrename.md
var batchRenameDocument string

const BatchRenameType = "batchrename"

func init() {
	addCmdDocumentInfo(BatchRenameType, batchRenameDocument)
}
