package docs

import _ "embed"

//go:embed sandbox_template_unpublish.md
var sandboxTemplateUnpublishDocument string

// SandboxTemplateUnpublishType is the document type for the sandbox template unpublish command.
const SandboxTemplateUnpublishType = "sandbox_template_unpublish"

func init() {
	addCmdDocumentInfo(SandboxTemplateUnpublishType, sandboxTemplateUnpublishDocument)
}
