package docs

import _ "embed"

//go:embed batchstat.md
var batchStatDocument string

const BatchStatType = "batchstat"

func init() {
	addCmdDocumentInfo(BatchStatType, batchStatDocument)
}
