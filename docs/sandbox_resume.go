package docs

import _ "embed"

//go:embed sandbox_resume.md
var sandboxResumeDocument string

// SandboxResumeType is the document type for the sandbox resume command.
const SandboxResumeType = "sandbox_resume"

func init() {
	addCmdDocumentInfo(SandboxResumeType, sandboxResumeDocument)
}
