package docs

import _ "embed"

//go:embed sandbox_logs.md
var sandboxLogsDocument string

// SandboxLogsType is the document type for the sandbox logs command.
const SandboxLogsType = "sandbox_logs"

func init() {
	addCmdDocumentInfo(SandboxLogsType, sandboxLogsDocument)
}
