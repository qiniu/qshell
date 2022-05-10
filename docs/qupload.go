package docs

import _ "embed"

//go:embed qupload.md
var qUploadDocument string

const QUploadType = "qupload"

func init() {
	addCmdDocumentInfo(QUploadType, qUploadDocument)
}
