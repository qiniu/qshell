package docs

import _ "embed"

//go:embed awsfetch.md
var AwsFetchDetailHelpString string
var AwsFetch = "awsfetch"

func init() {
	addCmdDocumentInfo(AwsFetch, AwsFetchDetailHelpString)
}
