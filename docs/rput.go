package docs

import _ "embed"

//go:embed rput.md
var rPutDocument string

const RPutType = "rput"

func init() {
	addCmdDocumentInfo(RPutType, rPutDocument)
}
