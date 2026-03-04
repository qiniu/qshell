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
