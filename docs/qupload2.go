package docs

import _ "embed"

//go:embed qupload2.md
var qUpload2Document string

const QUpload2Type = "qupload2"

func init() {
	addCmdDocumentInfo(QUpload2Type, qUpload2Document)
}
