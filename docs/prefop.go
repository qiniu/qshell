package docs

import _ "embed"

//go:embed prefop.md
var preFopDocument string

const PreFopType = "prefop"

func init() {
	addCmdDocumentInfo(PreFopType, preFopDocument)
}