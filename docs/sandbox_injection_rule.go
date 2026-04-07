package docs

import _ "embed"

//go:embed sandbox_injection_rule.md
var sandboxInjectionRuleDocument string

// SandboxInjectionRuleType is the document type for the sandbox injection-rule command.
const SandboxInjectionRuleType = "sandbox_injection_rule"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleType, sandboxInjectionRuleDocument)
}
