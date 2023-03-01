package docs

import (
	"github.com/qiniu/qshell/v2"
)

const QShellType = "qshell"

func init() {
	addCmdDocumentInfo(QShellType, qshell.ReadMeDocument)
}
