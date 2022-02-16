package docs

import _ "embed"

//go:embed buckets.md
var bucketsDocument string

const BucketsType = "buckets"

func init() {
	addCmdDocumentInfo(BucketsType, bucketsDocument)
}
