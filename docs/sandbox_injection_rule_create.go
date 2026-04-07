package docs

import _ "embed"

//go:embed sandbox_injection_rule_create.md
var sandboxInjectionRuleCreateDocument string

// SandboxInjectionRuleCreateType is the document type for the sandbox injection-rule create command.
const SandboxInjectionRuleCreateType = "sandbox_injection_rule_create"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleCreateType, sandboxInjectionRuleCreateDocument)
}
