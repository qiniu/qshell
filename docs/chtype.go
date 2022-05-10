package docs

import _ "embed"

//go:embed chtype.md
var changeTypeDocument string

const ChangeType = "chtype"

func init() {
	addCmdDocumentInfo(ChangeType, changeTypeDocument)
}
