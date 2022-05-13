package docs

import _ "embed"

//go:embed bucket.md
var bucketDocument string

const BucketType = "bucket"

func init() {
	addCmdDocumentInfo(BucketType, bucketDocument)
}

