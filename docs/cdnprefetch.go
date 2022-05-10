package docs

import _ "embed"

//go:embed cdnprefetch.md
var cdnPrefetchDocument string

const CdnPrefetchType = "cdnprefetch"

func init() {
	addCmdDocumentInfo(CdnPrefetchType, cdnPrefetchDocument)
}
