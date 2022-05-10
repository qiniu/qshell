package docs

import _ "embed"

//go:embed alilistbucket.md
var aliListBucketDocument string

const AliListBucket = "alilistbucket"

func init() {
	addCmdDocumentInfo(AliListBucket, aliListBucketDocument)
}
