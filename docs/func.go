package docs

import _ "embed"

//go:embed func.md
var funcDocument string

const FuncType = "func"

func init() {
	addCmdDocumentInfo(FuncType, funcDocument)
}
