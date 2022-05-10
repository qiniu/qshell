package docs

import _ "embed"

//go:embed listbucket.md
var listBucketDocument string

const ListBucketType = "listbucket"

func init() {
	addCmdDocumentInfo(ListBucketType, listBucketDocument)
}
