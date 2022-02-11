package docs

import _ "embed"

//go:embed account.md
var AccountDetailHelpString string
var Account = "account"

func init() {
	addCmdDocumentInfo(Account, AccountDetailHelpString)
}
