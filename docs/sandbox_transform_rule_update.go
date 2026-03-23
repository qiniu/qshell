package docs

import _ "embed"

//go:embed sandbox_transform_rule_update.md
var sandboxTransformRuleUpdateDocument string

// SandboxTransformRuleUpdateType is the document type for the sandbox transform-rule update command.
const SandboxTransformRuleUpdateType = "sandbox_transform_rule_update"

func init() {
	addCmdDocumentInfo(SandboxTransformRuleUpdateType, sandboxTransformRuleUpdateDocument)
}
