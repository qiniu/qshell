package docs

import _ "embed"

//go:embed stat.md
var statDocument string

const StatType = "stat"

func init() {
	addCmdDocumentInfo(StatType, statDocument)
}
