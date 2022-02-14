package docs

import _ "embed"

//go:embed b64encode.md
var B64EncodeDetailHelpString string
var B64Encode = "b64encode"

func init() {
	addCmdDocumentInfo(B64Encode, B64EncodeDetailHelpString)
}
