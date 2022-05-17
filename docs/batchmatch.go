package docs

import _ "embed"

//go:embed batchmatch.md
var batchMatchDocument string

const BatchMatchType = "batchmatch"

func init() {
	addCmdDocumentInfo(BatchMatchType, batchMatchDocument)
}
