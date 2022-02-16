package docs

import _ "embed"

//go:embed listbucket.md
var listBucketDocument string

const ListbucketType = "listbucket"

func init() {
	addCmdDocumentInfo(ListbucketType, listBucketDocument)
}
