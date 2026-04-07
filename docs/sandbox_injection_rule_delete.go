package docs

import _ "embed"

//go:embed sandbox_injection_rule_delete.md
var sandboxInjectionRuleDeleteDocument string

// SandboxInjectionRuleDeleteType is the document type for the sandbox injection-rule delete command.
const SandboxInjectionRuleDeleteType = "sandbox_injection_rule_delete"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleDeleteType, sandboxInjectionRuleDeleteDocument)
}
