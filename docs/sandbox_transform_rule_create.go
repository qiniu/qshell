package docs

import _ "embed"

//go:embed sandbox_transform_rule_create.md
var sandboxTransformRuleCreateDocument string

// SandboxTransformRuleCreateType is the document type for the sandbox transform-rule create command.
const SandboxTransformRuleCreateType = "sandbox_transform_rule_create"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleCreateType, sandboxTransformRuleCreateDocument)
}
