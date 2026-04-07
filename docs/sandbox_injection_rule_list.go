package docs

import _ "embed"

//go:embed sandbox_injection_rule_list.md
var sandboxInjectionRuleListDocument string

// SandboxInjectionRuleListType is the document type for the sandbox injection-rule list command.
const SandboxInjectionRuleListType = "sandbox_injection_rule_list"

func init() {
	addCmdDocumentInfo(SandboxInjectionRuleListType, sandboxInjectionRuleListDocument)
}
