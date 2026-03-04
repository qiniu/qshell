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

func TestSandboxLogsDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "logs")
}

func TestSandboxMetricsDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "metrics")
}

// === sandbox missing args tests ===
// These tests verify commands don't panic when required args are missing.
// Output is not captured (goes through cobra's Usage/fmt.Printf) but crash = test failure.

func TestSandboxConnectNoArgs(t *testing.T) {
	// connect without sandboxID should show usage, not panic
	test.RunCmdWithError("sandbox", "connect")
}

func TestSandboxLogsNoArgs(t *testing.T) {
	// logs without sandboxID should show usage, not panic
	test.RunCmdWithError("sandbox", "logs")
}

func TestSandboxMetricsNoArgs(t *testing.T) {
	// metrics without sandboxID should show usage, not panic
	test.RunCmdWithError("sandbox", "metrics")
}

func TestSandboxKillNoArgsNoAll(t *testing.T) {
	// kill without sandboxIDs and without --all should not panic
	// (will fail at API level with no API key, but should not crash)
	test.RunCmdWithError("sandbox", "kill")
}

func TestSandboxCreateNoArgs(t *testing.T) {
	// create without template should not panic
	// (operations.Create prints error and returns)
	test.RunCmdWithError("sandbox", "create")
}

// === sandbox --doc with flags tests ===
// Verify --doc still works when flags are also specified.
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

func TestSandboxListDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "list"}, "--state", "running")
}

func TestSandboxLogsDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "logs"}, "--level", "ERROR")
}

func TestSandboxMetricsDocumentWithFollow(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "metrics"}, "--follow")
}

func TestSandboxKillDocumentWithAll(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "kill"}, "--all")
}
