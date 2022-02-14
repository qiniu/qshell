package docs

import _ "embed"

//go:embed awslist.md
var AwsListDetailHelpString string
var AwsList = "awslist"

func init() {
	addCmdDocumentInfo(AwsList, AwsListDetailHelpString)
}
