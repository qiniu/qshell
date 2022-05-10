package docs

import _ "embed"

//go:embed privateurl.md
var privateUrlDocument string

const PrivateUrlType = "privateurl"

func init() {
	addCmdDocumentInfo(PrivateUrlType, privateUrlDocument)
}
