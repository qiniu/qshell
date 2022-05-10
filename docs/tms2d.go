package docs

import _ "embed"

//go:embed tms2d.md
var tms2dDocument string

const TMs2dType = "tms2d"

func init() {
	addCmdDocumentInfo(TMs2dType, tms2dDocument)
}
