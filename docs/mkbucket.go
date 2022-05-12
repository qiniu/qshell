package docs

import _ "embed"

//go:embed mkbucket.md
var mkBucketDocument string

const MkBucketDocument = "mkBucketDocument"

func init() {
	addCmdDocumentInfo(MkBucketDocument, mkBucketDocument)
}
