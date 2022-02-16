package docs

import _ "embed"

//go:embed m3u8delete.md
var m3u8DeleteDocument string

const M3u8DeleteType = "m3u8delete"

func init() {
	addCmdDocumentInfo(M3u8DeleteType, m3u8DeleteDocument)
}
