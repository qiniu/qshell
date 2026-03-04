package docs

import _ "embed"

//go:embed sandbox_template_builds.md
var sandboxTemplateBuildsDocument string

// SandboxTemplateBuildsType is the document type for the sandbox template builds command.
const SandboxTemplateBuildsType = "sandbox_template_builds"

func init() {
	addCmdDocumentInfo(SandboxTemplateBuildsType, sandboxTemplateBuildsDocument)
}
