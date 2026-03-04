package docs

import _ "embed"

//go:embed sandbox_template_build.md
var sandboxTemplateBuildDocument string

// SandboxTemplateBuildType is the document type for the sandbox template build command.
const SandboxTemplateBuildType = "sandbox_template_build"

func init() {
	addCmdDocumentInfo(SandboxTemplateBuildType, sandboxTemplateBuildDocument)
}
