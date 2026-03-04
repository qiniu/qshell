package docs

import _ "embed"

//go:embed sandbox_metrics.md
var sandboxMetricsDocument string

// SandboxMetricsType is the document type for the sandbox metrics command.
const SandboxMetricsType = "sandbox_metrics"

func init() {
	addCmdDocumentInfo(SandboxMetricsType, sandboxMetricsDocument)
}
