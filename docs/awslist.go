package docs

import _ "embed"

//go:embed awslist.md
var awsListDocument string

const AwsList = "awslist"

func init() {
	addCmdDocumentInfo(AwsList, awsListDocument)
}
