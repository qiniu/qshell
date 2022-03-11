package docs

import _ "embed"

//go:embed rename.md
var renameDocument string

const RenameType = "rename"

func init() {
	addCmdDocumentInfo(RenameType, renameDocument)
}
