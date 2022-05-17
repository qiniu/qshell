package docs

import _ "embed"

//go:embed match.md
var matchDocument string

const MatchType = "match"

func init() {
	addCmdDocumentInfo(MatchType, matchDocument)
}
