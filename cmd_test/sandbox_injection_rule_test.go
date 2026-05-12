//go:build unit

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func testInjectionRuleDocContains(t *testing.T, cmdParts []string, expected string, extraArgs ...string) {
	t.Helper()
	args := append(append([]string{}, cmdParts...), extraArgs...)
	args = append(args, test.DocumentOption)
	out, _ := test.RunCmdWithError(args...)
	if !strings.Contains(out, expected) {
		t.Fatalf("document output missing %q, got: %.120s", expected, out)
	}
}

func TestSandboxInjectionRuleDocument(t *testing.T) {
	testInjectionRuleDocContains(t, []string{"sandbox", "injection-rule"}, "`sandbox injection-rule`")
}

func TestSandboxInjectionRuleListDocument(t *testing.T) {
	testInjectionRuleDocContains(t, []string{"sandbox", "injection-rule", "list"}, "`sandbox injection-rule list`")
}

func TestSandboxInjectionRuleGetDocument(t *testing.T) {
	testInjectionRuleDocContains(t, []string{"sandbox", "injection-rule", "get"}, "`sandbox injection-rule get`", "rule-test")
}

func TestSandboxInjectionRuleCreateDocumentWithQiniu(t *testing.T) {
	testInjectionRuleDocContains(
		t,
		[]string{"sandbox", "injection-rule", "create"},
		"--type <openai|anthropic|gemini|qiniu|github|http>",
		"--name", "qiniu-default",
		"--type", "qiniu",
	)
}

func TestSandboxInjectionRuleUpdateDocumentWithQiniu(t *testing.T) {
	testInjectionRuleDocContains(
		t,
		[]string{"sandbox", "injection-rule", "update"},
		"`sandbox injection-rule update`",
		"rule-test",
		"--type", "qiniu",
	)
}

func TestSandboxInjectionRuleDeleteDocument(t *testing.T) {
	testInjectionRuleDocContains(t, []string{"sandbox", "injection-rule", "delete"}, "`sandbox injection-rule delete`", "rule-a")
}

func TestSandboxInjectionRuleDeleteNoArgs(t *testing.T) {
	// delete without ids and without --select should show usage and return safely
	test.RunCmdWithError("sandbox", "injection-rule", "delete")
}
