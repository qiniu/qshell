package docs

import _ "embed"

//go:embed urlencode.md
var urlEncodeDocument string

const UrlEncodeType = "urlencode"

func init() {
	addCmdDocumentInfo(UrlEncodeType, urlEncodeDocument)
}