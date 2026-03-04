package docs

import _ "embed"

//go:embed sandbox_template_get.md
var sandboxTemplateGetDocument string

// SandboxTemplateGetType is the document type for the sandbox template get command.
const SandboxTemplateGetType = "sandbox_template_get"

func init() {
	addCmdDocumentInfo(SandboxTemplateGetType, sandboxTemplateGetDocument)
}
