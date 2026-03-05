package operations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffold_Go(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold("my-app", "go", dir); err != nil {
		t.Fatalf("scaffold go failed: %v", err)
	}

	// Verify generated files exist
	for _, name := range []string{"main.go", "go.mod", "Makefile"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s failed: %v", name, err)
		}
		if len(data) == 0 {
			t.Errorf("%s is empty", name)
		}
	}

	// Verify template substitution in main.go
	mainGo, _ := os.ReadFile(filepath.Join(dir, "main.go"))
	if !strings.Contains(string(mainGo), "my-app") {
		t.Error("main.go should contain template name 'my-app'")
	}

	// Verify go.mod contains module name
	goMod, _ := os.ReadFile(filepath.Join(dir, "go.mod"))
	if !strings.Contains(string(goMod), "module my-app") {
		t.Error("go.mod should contain 'module my-app'")
	}
}

func TestScaffold_TypeScript(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold("ts-project", "typescript", dir); err != nil {
		t.Fatalf("scaffold typescript failed: %v", err)
	}

	for _, name := range []string{"template.ts", "package.json"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s failed: %v", name, err)
		}
		if len(data) == 0 {
			t.Errorf("%s is empty", name)
		}
	}

	// Verify package.json contains project name
	pkgJSON, _ := os.ReadFile(filepath.Join(dir, "package.json"))
	if !strings.Contains(string(pkgJSON), "ts-project") {
		t.Error("package.json should contain project name 'ts-project'")
	}
}

func TestScaffold_Python(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold("py-script", "python", dir); err != nil {
		t.Fatalf("scaffold python failed: %v", err)
	}

	for _, name := range []string{"template.py", "requirements.txt"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s failed: %v", name, err)
		}
		if len(data) == 0 {
			t.Errorf("%s is empty", name)
		}
	}

	// Verify template.py contains project name
	tmplPy, _ := os.ReadFile(filepath.Join(dir, "template.py"))
	if !strings.Contains(string(tmplPy), "py-script") {
		t.Error("template.py should contain project name 'py-script'")
	}
}

func TestScaffold_UnsupportedLanguage(t *testing.T) {
	dir := t.TempDir()
	err := scaffold("test", "rust", dir)
	if err == nil {
		t.Fatal("scaffold should fail for unsupported language")
	}
	if !strings.Contains(err.Error(), "unsupported language") {
		t.Errorf("error = %v, want 'unsupported language'", err)
	}
}

func TestScaffold_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	if err := scaffold("test-proj", "go", dir); err != nil {
		t.Fatalf("scaffold with nested dir failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "main.go")); os.IsNotExist(err) {
		t.Error("main.go should exist in created directory")
	}
}
