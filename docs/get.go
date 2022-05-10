package docs

import _ "embed"

//go:embed get.md
var getDocument string

const GetType = "get"

func init() {
	addCmdDocumentInfo(GetType, getDocument)
}
