package docs

import _ "embed"

//go:embed unzip.md
var unzipDocument string

const UnzipType = "unzip"

func init() {
	addCmdDocumentInfo(UnzipType, unzipDocument)
}
