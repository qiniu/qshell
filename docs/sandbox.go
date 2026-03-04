package docs

import _ "embed"

//go:embed sandbox.md
var sandboxDocument string

// SandboxType is the document type for the sandbox command.
const SandboxType = "sandbox"

func init() {
	addCmdDocumentInfo(SandboxType, sandboxDocument)
}
