package docs

import _ "embed"

//go:embed tns2d.md
var tns2dDocument string

const TNs2dType = "tns2d"

func init() {
	addCmdDocumentInfo(TNs2dType, tns2dDocument)
}
