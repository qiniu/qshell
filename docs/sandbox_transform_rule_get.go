package docs

import _ "embed"

//go:embed sandbox_transform_rule_get.md
var sandboxTransformRuleGetDocument string

// SandboxTransformRuleGetType is the document type for the sandbox transform-rule get command.
const SandboxTransformRuleGetType = "sandbox_transform_rule_get"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleGetType, sandboxTransformRuleGetDocument)
}
