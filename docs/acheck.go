package docs

import _ "embed"

//go:embed acheck.md
var aCheckDocument string

const ACheckType = "acheck"

func init() {
	addCmdDocumentInfo(ACheckType, aCheckDocument)
}
