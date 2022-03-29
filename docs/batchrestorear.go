package docs

import _ "embed"

//go:embed batchrestorear.md
var batchRestoreArchiveDocument string

const BatchRestoreArchiveType = "batchrestorear"

func init() {
	addCmdDocumentInfo(BatchRestoreArchiveType, batchRestoreArchiveDocument)
}
