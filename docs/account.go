package docs

import _ "embed"

//go:embed account.md
var AccountDetailHelpString string
var Account = "abfetch"

func init() {
	addCmdDocumentInfo(Account, AccountDetailHelpString)
}
