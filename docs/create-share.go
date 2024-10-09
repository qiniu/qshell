package docs

import _ "embed"

//go:embed create-share.md
var createShareDocument string

const CreateShareType = "create-share"

func init() {
	addCmdDocumentInfo(CreateShareType, createShareDocument)
}
