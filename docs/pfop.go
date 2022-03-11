package docs

import _ "embed"

//go:embed pfop.md
var pFopDocument string

const PFopType = "pfop"

func init() {
	addCmdDocumentInfo(PFopType, pFopDocument)
}
