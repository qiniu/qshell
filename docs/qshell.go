package docs

import _ "embed"

//go:embed qshell.md
var qshellDocument string

const QShellType = "qshell"

func init() {
	addCmdDocumentInfo(QShellType, qshellDocument)
}
