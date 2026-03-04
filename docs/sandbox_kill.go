package docs

import _ "embed"

//go:embed sandbox_kill.md
var sandboxKillDocument string

// SandboxKillType is the document type for the sandbox kill command.
const SandboxKillType = "sandbox_kill"

func init() {
	addCmdDocumentInfo(SandboxKillType, sandboxKillDocument)
}
