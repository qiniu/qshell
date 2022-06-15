package docs

import _ "embed"

//go:embed forbidden.md
var forbiddenDocument string

const ForbiddenType = "forbidden"

func init() {
	addCmdDocumentInfo(ForbiddenType, forbiddenDocument)
}
