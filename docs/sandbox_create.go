package docs

import _ "embed"

//go:embed sandbox_create.md
var sandboxCreateDocument string

// SandboxCreateType is the document type for the sandbox create command.
const SandboxCreateType = "sandbox_create"

func init() {
	addCmdDocumentInfo(SandboxCreateType, sandboxCreateDocument)
}
