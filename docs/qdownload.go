package docs

import _ "embed"

//go:embed qdownload.md
var qDownloadDocument string

const QDownloadType = "qdownload"

func init() {
	addCmdDocumentInfo(QDownloadType, qDownloadDocument)
}