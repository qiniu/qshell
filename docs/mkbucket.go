package docs

import _ "embed"

//go:embed mkbucket.md
var mkBucketDocument string

const MkBucketType = "mkBucketDocument"

func init() {
	addCmdDocumentInfo(MkBucketType, mkBucketDocument)
}
