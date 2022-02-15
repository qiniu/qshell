package docs

import _ "embed"

//go:embed account.md
var accountDocument string

const Account = "account"

func init() {
	addCmdDocumentInfo(Account, accountDocument)
}
