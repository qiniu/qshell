package docs

import _ "embed"

//go:embed batchchgm.md
var batchChangeMimeTypeDocument string

const BatchChangeMimeType = "batchchgm"

func init() {
	addCmdDocumentInfo(BatchChangeMimeType, batchChangeMimeTypeDocument)
}
