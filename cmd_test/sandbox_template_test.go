//go:build unit

package cmd

import (
	"testing"
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
