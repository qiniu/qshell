package docs

import _ "embed"

//go:embed ts2d.md
var ts2dDocument string

const TS2dType = "ts2d"

func init() {
	addCmdDocumentInfo(TS2dType, ts2dDocument)
}
