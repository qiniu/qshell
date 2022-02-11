package docs

import _ "embed"

//go:embed user.md
var UserDetailHelpString string
var User = "user"

func init() {
	addCmdDocumentInfo(User, UserDetailHelpString)
}
