package docs

import _ "embed"

//go:embed d2ts.md
var dateToTimestampDocument string

const DateToTimestamp = "d2ts"

func init() {
	addCmdDocumentInfo(DateToTimestamp, dateToTimestampDocument)
}
