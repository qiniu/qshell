package docs

import _ "embed"

//go:embed sandbox_transform_rule_list.md
var sandboxTransformRuleListDocument string

// SandboxTransformRuleListType is the document type for the sandbox transform-rule list command.
const SandboxTransformRuleListType = "sandbox_transform_rule_list"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleListType, sandboxTransformRuleListDocument)
}
