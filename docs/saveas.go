package docs

import _ "embed"

//go:embed saveas.md
var saveAsDocument string

const SaveAsType = "saveas"

func init() {
	addCmdDocumentInfo(SaveAsType, saveAsDocument)
}
