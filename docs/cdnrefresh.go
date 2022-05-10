package docs

import _ "embed"

//go:embed cdnrefresh.md
var cdnRefreshDocument string

const CdnRefreshType = "cdnrefresh"

func init() {
	addCmdDocumentInfo(CdnRefreshType, cdnRefreshDocument)
}
