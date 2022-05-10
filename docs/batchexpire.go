package docs

import _ "embed"

//go:embed batchexpire.md
var batchExpireDocument string

const BatchExpireType = "batchexpire"

func init() {
	addCmdDocumentInfo(BatchExpireType, batchExpireDocument)
}
