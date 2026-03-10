package docs

import _ "embed"

//go:embed sandbox_pause.md
var sandboxPauseDocument string

// SandboxPauseType is the document type for the sandbox pause command.
const SandboxPauseType = "sandbox_pause"

func init() {
	addCmdDocumentInfo(SandboxPauseType, sandboxPauseDocument)
}
