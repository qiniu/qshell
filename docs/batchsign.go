package docs

import _ "embed"

//go:embed batchsign.md
var batchSignDocument string

const BatchSignType = "batchsign"

func init() {
	addCmdDocumentInfo(BatchSignType, batchSignDocument)
}
