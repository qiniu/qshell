package docs

import _ "embed"

//go:embed mirrorupdate.md
var mirrorUpdateDocument string

const MirrorUpdateType = "mirrorupdate"

func init() {
	addCmdDocumentInfo(MirrorUpdateType, mirrorUpdateDocument)
}
