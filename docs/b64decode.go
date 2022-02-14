package docs

import _ "embed"

//go:embed b64decode.md
var B64DecodeDetailHelpString string
var B64Decode = "b64decode"

func init() {
	addCmdDocumentInfo(B64Decode, B64DecodeDetailHelpString)
}
