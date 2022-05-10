package docs

import _ "embed"

//go:embed ip.md
var ipDocument string

const IPType = "ip"

func init() {
	addCmdDocumentInfo(IPType, ipDocument)
}
