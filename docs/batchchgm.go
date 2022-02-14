package docs

import _ "embed"

//go:embed batchchgm.md
var BatchChangeMimeTypeDetailHelpString string
var BatchChangeMimeType = "batchchgm"

func init() {
	addCmdDocumentInfo(BatchChangeMimeType, BatchChangeMimeTypeDetailHelpString)
}

