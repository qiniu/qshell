package docs

import _ "embed"

//go:embed sandbox_injection_rule_update.md
var sandboxInjectionRuleUpdateDocument string

// SandboxInjectionRuleUpdateType is the document type for the sandbox injection-rule update command.
const SandboxInjectionRuleUpdateType = "sandbox_injection_rule_update"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleUpdateType, sandboxInjectionRuleUpdateDocument)
}
