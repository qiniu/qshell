//go:build unit

package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

// testSubcommandDocument tests --doc output for a subcommand with multiple args.
// Unlike test.TestDocument which only works for single-word commands,
// this function splits the command path into separate args.
func testSubcommandDocument(t *testing.T, cmdParts ...string) {
	t.Helper()
	cmdName := strings.Join(cmdParts, " ")
	prefix := fmt.Sprintf("# 简介\n`%s`", cmdName)
	args := append(cmdParts, test.DocumentOption)
	result, _ := test.RunCmdWithError(args...)
	if !strings.HasPrefix(result, prefix) {
		t.Fatalf("document test fail for cmd: %s, got prefix: %.80s", cmdName, result)
	}
}

// testSubcommandDocumentWithFlags passes extra flags but only checks the base command name in prefix.
func testSubcommandDocumentWithFlags(t *testing.T, cmdParts []string, extraFlags ...string) {
	t.Helper()
	cmdName := strings.Join(cmdParts, " ")
	prefix := fmt.Sprintf("# 简介\n`%s`", cmdName)
	args := append(cmdParts, extraFlags...)
	args = append(args, test.DocumentOption)
	result, _ := test.RunCmdWithError(args...)
	if !strings.HasPrefix(result, prefix) {
		t.Fatalf("document test fail for cmd: %s with flags %v, got prefix: %.80s", cmdName, extraFlags, result)
	}
}

// === sandbox document tests ===

func TestSandboxDocument(t *testing.T) {
	test.TestDocument("sandbox", t)
}

func TestSandboxListDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "list")
}

func TestSandboxCreateDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "create")
}

func TestSandboxConnectDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "connect")
}

func TestSandboxKillDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "kill")
}

func TestSandboxPauseDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "pause")
}

func TestSandboxResumeDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "resume")
}

func TestSandboxExecDocument(t *testing.T) {
	// exec uses cobra.MinimumNArgs(1), so we must pass a sandboxID for --doc to work
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "exec"}, "sb-test")
}

func TestSandboxLogsDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "logs")
}

func TestSandboxMetricsDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "metrics")
}

// === sandbox missing args tests ===
// These tests verify commands don't panic when required args are missing.

func TestSandboxConnectNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "connect")
}

func TestSandboxLogsNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "logs")
}

func TestSandboxMetricsNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "metrics")
}

func TestSandboxKillNoArgsNoAll(t *testing.T) {
	test.RunCmdWithError("sandbox", "kill")
}

func TestSandboxCreateNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "create")
}

func TestSandboxPauseNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "pause")
}

func TestSandboxResumeNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "resume")
}

// === sandbox --doc with flags tests ===

func TestSandboxListDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "list"}, "--state", "running")
}

func TestSandboxCreateDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "create"}, "--timeout", "300", "--detach", "--auto-pause")
}

func TestSandboxCreateDocumentWithEnvVar(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "create"}, "-e", "FOO=bar", "-e", "BAZ=qux")
}

func TestSandboxKillDocumentWithAll(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "kill"}, "--all")
}

func TestSandboxPauseDocumentWithAll(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "pause"}, "--all")
}

func TestSandboxPauseDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "pause"}, "--all", "-s", "running", "-m", "env=dev")
}

func TestSandboxResumeDocumentWithAll(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "resume"}, "--all")
}

func TestSandboxResumeDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "resume"}, "--all", "-m", "env=staging")
}

func TestSandboxExecDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "exec"}, "sb-test", "-b", "-c", "/app", "-u", "root")
}

func TestSandboxLogsDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "logs"}, "--level", "ERROR")
}

func TestSandboxMetricsDocumentWithFollow(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "metrics"}, "--follow")
}
