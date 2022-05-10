package docs

import _ "embed"

//go:embed reqid.md
var reqIdDocument string

const ReqIdType = "reqid"

func init() {
	addCmdDocumentInfo(ReqIdType, reqIdDocument)
}
