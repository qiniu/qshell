package docs

import _ "embed"

//go:embed sandbox_injection_rule_get.md
var sandboxInjectionRuleGetDocument string

// SandboxInjectionRuleGetType is the document type for the sandbox injection-rule get command.
const SandboxInjectionRuleGetType = "sandbox_injection_rule_get"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleGetType, sandboxInjectionRuleGetDocument)
}
