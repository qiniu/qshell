package docs

import _ "embed"

//go:embed sandbox_template_publish.md
var sandboxTemplatePublishDocument string

// SandboxTemplatePublishType is the document type for the sandbox template publish command.
const SandboxTemplatePublishType = "sandbox_template_publish"

func init() {
	addCmdDocumentInfo(SandboxTemplatePublishType, sandboxTemplatePublishDocument)
}
