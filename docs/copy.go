package docs

import _ "embed"

//go:embed copy.md
var copyDocument string

const CopyType = "copy"

func init() {
	addCmdDocumentInfo(CopyType, copyDocument)
}
