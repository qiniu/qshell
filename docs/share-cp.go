package docs

import _ "embed"

//go:embed share-cp.md
var shareCpDocument string

const ShareCpType = "share-cp"

func init() {
	addCmdDocumentInfo(ShareCpType, shareCpDocument)
}
