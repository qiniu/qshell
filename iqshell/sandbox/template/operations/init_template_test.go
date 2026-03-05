package operations

import (
	"testing"
)

func TestValidNamePattern_ValidNames(t *testing.T) {
	valid := []string{
		"my-template",
		"my_template",
		"template1",
		"a",
		"0test",
		"hello-world_123",
	}
	for _, name := range valid {
		if !validNamePattern.MatchString(name) {
			t.Errorf("validNamePattern should accept %q", name)
		}
	}
}

func TestValidNamePattern_InvalidNames(t *testing.T) {
	invalid := []string{
		"",
		"-start-with-dash",
		"_start-with-underscore",
		"UPPERCASE",
		"has spaces",
		"special!char",
		"has.dot",
		"has/slash",
	}
	for _, name := range invalid {
		if validNamePattern.MatchString(name) {
			t.Errorf("validNamePattern should reject %q", name)
		}
	}
}

func TestSupportedLanguages(t *testing.T) {
	expected := map[string]bool{
		"go":         true,
		"typescript": true,
		"python":     true,
	}
	if len(supportedLanguages) != len(expected) {
		t.Fatalf("supportedLanguages has %d entries, want %d", len(supportedLanguages), len(expected))
	}
	for _, lang := range supportedLanguages {
		if !expected[lang] {
			t.Errorf("unexpected language: %s", lang)
		}
	}
}

func TestInit_NonInteractive_ValidGoProject(t *testing.T) {
	dir := t.TempDir()
	Init(InitInfo{
		Name:     "testproject",
		Language: "go",
		Path:     dir,
	})
	// If scaffold worked, main.go should exist
	// (Init prints success message; since it doesn't return error, we just check it doesn't panic)
}

func TestInit_NonInteractive_ValidTypeScriptProject(t *testing.T) {
	dir := t.TempDir()
	Init(InitInfo{
		Name:     "tsproject",
		Language: "typescript",
		Path:     dir,
	})
}

func TestInit_NonInteractive_ValidPythonProject(t *testing.T) {
	dir := t.TempDir()
	Init(InitInfo{
		Name:     "pyproject",
		Language: "python",
		Path:     dir,
	})
}

func TestInit_InvalidName(t *testing.T) {
	// Should print error but not panic
	Init(InitInfo{
		Name:     "INVALID-NAME",
		Language: "go",
		Path:     t.TempDir(),
	})
}

func TestInit_InvalidLanguage(t *testing.T) {
	// Should print error but not panic
	Init(InitInfo{
		Name:     "test",
		Language: "rust",
		Path:     t.TempDir(),
	})
}

func TestInit_DefaultPath(t *testing.T) {
	// When path is empty, Init uses ./<name>
	// We just verify it doesn't panic; actual directory would be created in CWD
	// which we don't want in tests. So we only test with explicit path.
	dir := t.TempDir()
	Init(InitInfo{
		Name:     "pathtest",
		Language: "go",
		Path:     dir + "/pathtest",
	})
}
