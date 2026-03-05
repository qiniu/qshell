package docs

import _ "embed"

//go:embed sandbox_template_init.md
var sandboxTemplateInitDocument string

// SandboxTemplateInitType is the document type for the sandbox template init command.
const SandboxTemplateInitType = "sandbox_template_init"

func init() {
	addCmdDocumentInfo(SandboxTemplateInitType, sandboxTemplateInitDocument)
}
