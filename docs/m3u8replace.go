package docs

import _ "embed"

//go:embed m3u8replace.md
var m3u8ReplaceDocument string

const M3u8ReplaceType = "m3u8replace"

func init() {
	addCmdDocumentInfo(M3u8ReplaceType, m3u8ReplaceDocument)
}
