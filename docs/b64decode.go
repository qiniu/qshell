package docs

import _ "embed"

//go:embed b64decode.md
var b64DecodeDocument string

const B64Decode = "b64decode"

func init() {
	addCmdDocumentInfo(B64Decode, b64DecodeDocument)
}
