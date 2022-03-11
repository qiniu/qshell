package docs

import _ "embed"

//go:embed qetag.md
var qTagDocument string

const QTagType = "qetag"

func init() {
	addCmdDocumentInfo(QTagType, qTagDocument)
}
