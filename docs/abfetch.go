package docs

import _ "embed"

//go:embed abfetch.md
var ABFetchDetailHelpString string
var ABFetch = "abfetch"

func init() {
	addCmdDocumentInfo(ABFetch, ABFetchDetailHelpString)
}