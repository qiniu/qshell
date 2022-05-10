package docs

import _ "embed"

//go:embed move.md
var moveDocument string

const MoveType = "move"

func init() {
	addCmdDocumentInfo(MoveType, moveDocument)
}
