package docs

import _ "embed"

//go:embed domains.md
var domainsDocument string

const DomainsType = "domains"

func init() {
	addCmdDocumentInfo(DomainsType, domainsDocument)
}