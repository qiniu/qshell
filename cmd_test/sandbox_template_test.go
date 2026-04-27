//go:build unit

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

// === sandbox template document tests ===

func TestSandboxTemplateDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template")
}

func TestSandboxTemplateListDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "list")
}

func TestSandboxTemplateGetDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "get")
}

func TestSandboxTemplateDeleteDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "delete")
}

func TestSandboxTemplateBuildDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "build")
}

func TestSandboxTemplateBuildsDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "builds")
}

func TestSandboxTemplatePublishDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "publish")
}

func TestSandboxTemplateUnpublishDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "unpublish")
}

func TestSandboxTemplateInitDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "init")
}

func TestSandboxTemplateConfigDocument(t *testing.T) {
	testSubcommandDocument(t, "sandbox", "template", "config")
}

// === sandbox template missing args tests ===

func TestSandboxTemplateDeleteNoArgs(t *testing.T) {
	// delete without templateIDs and without --select should show usage, not panic
	test.RunCmdWithError("sandbox", "template", "delete")
}

func TestSandboxTemplatePublishNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "template", "publish")
}

func TestSandboxTemplateUnpublishNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "template", "unpublish")
}

func TestSandboxTemplateGetNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "template", "get")
}

func TestSandboxTemplateBuildNoArgs(t *testing.T) {
	// build without --name or --template-id should print error, not panic
	test.RunCmdWithError("sandbox", "template", "build")
}

func TestSandboxTemplateBuildsNoArgs(t *testing.T) {
	test.RunCmdWithError("sandbox", "template", "builds")
}

func TestSandboxTemplateInitNoArgs(t *testing.T) {
	// init without args triggers interactive prompt; in non-TTY test mode
	// huh form fails immediately, which is handled gracefully
	test.RunCmdWithError("sandbox", "template", "init")
}

// === sandbox template --doc with flags tests ===

func TestSandboxTemplateListDocumentWithFormat(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "list"}, "--format", "json")
}

func TestSandboxTemplateDeleteDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "delete"}, "--yes")
}

func TestSandboxTemplatePublishDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "publish"}, "--yes")
}

func TestSandboxTemplateUnpublishDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "unpublish"}, "--yes")
}

func TestSandboxTemplateBuildDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "build"}, "--no-cache")
}

func TestSandboxTemplateBuildDocumentWithDockerfile(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "build"}, "--dockerfile", "./Dockerfile")
}

func TestSandboxTemplateBuildDocumentWithDockerfileAndPath(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "build"}, "--dockerfile", "./Dockerfile", "--path", "./context")
}

func TestSandboxTemplateInitDocumentWithFlags(t *testing.T) {
	testSubcommandDocumentWithFlags(t, []string{"sandbox", "template", "init"}, "--name", "test")
}

func TestSandboxTemplateConfigDocumentWithAlias(t *testing.T) {
	result, _ := test.RunCmdWithError("sbx", "tpl", "cfg", test.DocumentOption)
	if !strings.Contains(result, "`sandbox template config`") {
		t.Fatalf("document test fail for sbx tpl cfg, got prefix: %.80s", result)
	}
}
