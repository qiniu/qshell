package docs

import _ "embed"

//go:embed sync.md
var syncDocument string

const SyncType = "sync"

func init() {
	addCmdDocumentInfo(SyncType, syncDocument)
}