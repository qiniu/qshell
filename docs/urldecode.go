package docs

import _ "embed"

//go:embed urldecode.md
var urlDecodeDocument string

const UrlDecodeType = "urldecode"

func init() {
	addCmdDocumentInfo(UrlDecodeType, urlDecodeDocument)
}
