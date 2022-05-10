package docs

import _ "embed"

//go:embed fput.md
var formPutDocument string

const FormPutType = "fput"

func init() {
	addCmdDocumentInfo(FormPutType, formPutDocument)
}
