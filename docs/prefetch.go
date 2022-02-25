package docs

import _ "embed"

//go:embed prefetch.md
var prefetchDocument string

const PrefetchType = "prefetch"

func init() {
	addCmdDocumentInfo(PrefetchType, prefetchDocument)
}
