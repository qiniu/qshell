package docs

import _ "embed"

//go:embed d2ts.md
var dateToTimestampDocument string

const DateToTimestampType = "d2ts"

func init() {
	addCmdDocumentInfo(DateToTimestampType, dateToTimestampDocument)
}
