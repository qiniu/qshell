package docs

import _ "embed"

//go:embed batchchtype.md
var batchChangeTypeDocument string

const BatchChangeType = "batchchtype"

func init() {
	addCmdDocumentInfo(BatchChangeType, batchChangeTypeDocument)
}
