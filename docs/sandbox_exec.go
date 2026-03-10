package docs

import _ "embed"

//go:embed sandbox_exec.md
var sandboxExecDocument string

// SandboxExecType is the document type for the sandbox exec command.
const SandboxExecType = "sandbox_exec"

func init() {
	addCmdDocumentInfo(SandboxExecType, sandboxExecDocument)
}
