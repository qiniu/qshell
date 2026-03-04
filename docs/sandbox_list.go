package docs

import _ "embed"

//go:embed sandbox_list.md
var sandboxListDocument string

// SandboxListType is the document type for the sandbox list command.
const SandboxListType = "sandbox_list"

func init() {
	addCmdDocumentInfo(SandboxListType, sandboxListDocument)
}
