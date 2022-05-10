package docs

import _ "embed"

//go:embed restorear.md
var restoreArchiveDocument string

const RestoreArchiveType = "restorear"

func init() {
	addCmdDocumentInfo(RestoreArchiveType, restoreArchiveDocument)
}
