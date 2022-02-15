package docs

import _ "embed"

//go:embed batchdelete.md
var batchDeleteDocument string

const BatchDeleteType = "batchdelete"

func init() {
	addCmdDocumentInfo(BatchDeleteType, batchDeleteDocument)
}
