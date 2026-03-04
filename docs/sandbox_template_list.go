package docs

import _ "embed"

//go:embed sandbox_template_list.md
var sandboxTemplateListDocument string

// SandboxTemplateListType is the document type for the sandbox template list command.
const SandboxTemplateListType = "sandbox_template_list"

func init() {
	addCmdDocumentInfo(SandboxTemplateListType, sandboxTemplateListDocument)
}
