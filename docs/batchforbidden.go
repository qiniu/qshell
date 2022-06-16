package docs

import _ "embed"

//go:embed batchforbidden.md
var batchForbiddenDocument string

const BatchForbiddenType = "batchforbidden"

func init() {
	addCmdDocumentInfo(BatchForbiddenType, batchForbiddenDocument)
}
