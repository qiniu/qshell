package docs

import _ "embed"

//go:embed delete.md
var deleteDocument string

const DeleteType = "delete"

func init() {
	addCmdDocumentInfo(DeleteType, deleteDocument)
}