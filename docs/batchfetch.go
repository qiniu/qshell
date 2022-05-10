package docs

import _ "embed"

//go:embed batchfetch.md
var batchFetchDocument string

const BatchFetchType = "batchfetch"

func init() {
	addCmdDocumentInfo(BatchFetchType, batchFetchDocument)
}
