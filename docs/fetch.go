package docs

import _ "embed"

//go:embed fetch.md
var fetchDocument string

const FetchType = "fetch"

func init() {
	addCmdDocumentInfo(FetchType, fetchDocument)
}