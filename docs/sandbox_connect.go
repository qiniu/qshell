package docs

import _ "embed"

//go:embed sandbox_connect.md
var sandboxConnectDocument string

// SandboxConnectType is the document type for the sandbox connect command.
const SandboxConnectType = "sandbox_connect"

func init() {
	addCmdDocumentInfo(SandboxConnectType, sandboxConnectDocument)
}
