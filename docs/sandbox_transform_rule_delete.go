package docs

import _ "embed"

//go:embed sandbox_transform_rule_delete.md
var sandboxTransformRuleDeleteDocument string

// SandboxTransformRuleDeleteType is the document type for the sandbox transform-rule delete command.
const SandboxTransformRuleDeleteType = "sandbox_transform_rule_delete"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleDeleteType, sandboxTransformRuleDeleteDocument)
}
