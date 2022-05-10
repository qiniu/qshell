package docs

import _ "embed"

//go:embed chgm.md
var changeMimeDocument string

const ChangeMimeType = "chgm"

func init() {
	addCmdDocumentInfo(ChangeMimeType, changeMimeDocument)
}
