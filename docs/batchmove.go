package docs

import _ "embed"

//go:embed batchmove.md
var batchMoveDocument string

const BatchMoveType = "batchmove"

func init() {
	addCmdDocumentInfo(BatchMoveType, batchMoveDocument)
}
