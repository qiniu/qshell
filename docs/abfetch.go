package docs

import _ "embed"

//go:embed abfetch.md
var abFetchDocument string

const ABFetch = "abfetch"

func init() {
	addCmdDocumentInfo(ABFetch, abFetchDocument)
}
