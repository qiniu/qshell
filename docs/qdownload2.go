package docs

import _ "embed"

//go:embed qdownload2.md
var qDownload2Document string

const QDownload2Type = "qdownload2"

func init() {
	addCmdDocumentInfo(QDownload2Type, qDownload2Document)
}
