package docs

import _ "embed"

//go:embed b64encode.md
var b64EncodeDocument string

const B64Encode = "b64encode"

func init() {
	addCmdDocumentInfo(B64Encode, b64EncodeDocument)
}
