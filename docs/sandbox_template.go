package docs

import _ "embed"

//go:embed sandbox_template.md
var sandboxTemplateDocument string

// SandboxTemplateType is the document type for the sandbox template command.
const SandboxTemplateType = "sandbox_template"

func init() {
	addCmdDocumentInfo(SandboxTemplateType, sandboxTemplateDocument)
}
