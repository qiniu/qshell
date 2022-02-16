package docs

import _ "embed"

//go:embed expire.md
var expireDocument string

const ExpireType = "expire"

func init() {
	addCmdDocumentInfo(ExpireType, expireDocument)
}