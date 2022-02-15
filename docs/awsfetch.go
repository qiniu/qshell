package docs

import _ "embed"

//go:embed awsfetch.md
var awsFetchDocument string

const AwsFetch = "awsfetch"

func init() {
	addCmdDocumentInfo(AwsFetch, awsFetchDocument)
}
