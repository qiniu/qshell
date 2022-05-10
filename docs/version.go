package docs

import _ "embed"

//go:embed version.md
var versionDocument string

const VersionType = "version"

func init() {
	addCmdDocumentInfo(VersionType, versionDocument)
}
