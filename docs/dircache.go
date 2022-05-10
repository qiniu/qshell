package docs

import _ "embed"

//go:embed dircache.md
var dirCacheDocument string

const DirCacheType = "dircache"

func init() {
	addCmdDocumentInfo(DirCacheType, dirCacheDocument)
}
