package docs

import _ "embed"

//go:embed sandbox_template_delete.md
var sandboxTemplateDeleteDocument string

// SandboxTemplateDeleteType is the document type for the sandbox template delete command.
const SandboxTemplateDeleteType = "sandbox_template_delete"

func init() {
	addCmdDocumentInfo(SandboxTemplateDeleteType, sandboxTemplateDeleteDocument)
}
