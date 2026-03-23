package docs

import _ "embed"

//go:embed sandbox_transform_rule.md
var sandboxTransformRuleDocument string

// SandboxTransformRuleType is the document type for the sandbox transform-rule command.
const SandboxTransformRuleType = "sandbox_transform_rule"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleType, sandboxTransformRuleDocument)
}
