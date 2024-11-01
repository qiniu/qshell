package docs

import _ "embed"

//go:embed share-ls.md
var shareLsDocument string

const ShareLsType = "share-ls"

func init() {
	addCmdDocumentInfo(ShareLsType, shareLsDocument)
}
