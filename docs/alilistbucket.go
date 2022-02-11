package docs

import _ "embed"

//go:embed alilistbucket.md
var AliListBucketDetailHelpString string
var AliListBucket = "alilistbucket"

func init() {
	addCmdDocumentInfo(AliListBucket, AliListBucketDetailHelpString)
}
