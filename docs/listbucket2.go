package docs

import _ "embed"

//go:embed listbucket2.md
var listBucket2Document string

const ListBucket2Type = "listbucket2"

func init() {
	addCmdDocumentInfo(ListBucket2Type, listBucket2Document)
}
