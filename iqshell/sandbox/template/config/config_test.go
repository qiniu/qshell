//go:build unit

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(filepath.Join(dir, "nope.toml"))
	assert.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	content := `template_id = "tmpl-abc"
name = "demo"
dockerfile = "./Dockerfile"
cpu_count = 2
memory_mb = 2048
no_cache = true
`
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))

	cfg, err := Load(p)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "tmpl-abc", cfg.TemplateID)
	assert.Equal(t, "demo", cfg.Name)
	assert.Equal(t, "./Dockerfile", cfg.Dockerfile)
	assert.Equal(t, int32(2), cfg.CPUCount)
	assert.Equal(t, int32(2048), cfg.MemoryMB)
	assert.True(t, cfg.NoCache)
	assert.Equal(t, p, cfg.SourcePath())
}

func TestLoad_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(p, []byte("this is = not = valid"), 0644))

	cfg, err := Load(p)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestFindInCwd_Found(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, DefaultFileName)
	require.NoError(t, os.WriteFile(p, []byte(""), 0644))

	got, err := FindInDir(dir)
	require.NoError(t, err)
	assert.Equal(t, p, got)
}

func TestFindInCwd_NotFound(t *testing.T) {
	dir := t.TempDir()
	got, err := FindInDir(dir)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestApplyTo_FileOnly(t *testing.T) {
	cfg := &FileConfig{
		TemplateID: "tmpl-x",
		Name:       "demo",
		Dockerfile: "./Dockerfile",
		CPUCount:   2,
		MemoryMB:   2048,
		NoCache:    true,
		defined:    map[string]bool{"no_cache": true},
	}
	dst := BuildFields{}
	overrides := cfg.ApplyTo(&dst)
	assert.Equal(t, "tmpl-x", dst.TemplateID)
	assert.Equal(t, "demo", dst.Name)
	assert.Equal(t, "./Dockerfile", dst.Dockerfile)
	assert.Equal(t, int32(2), dst.CPUCount)
	assert.Equal(t, int32(2048), dst.MemoryMB)
	assert.True(t, dst.NoCache)
	assert.Empty(t, overrides)
}

func TestApplyTo_CLIOverride(t *testing.T) {
	cfg := &FileConfig{
		Name:     "from-file",
		CPUCount: 2,
		NoCache:  false,
	}
	dst := BuildFields{
		Name:     "from-cli",
		CPUCount: 4,
		NoCache:  true,
	}
	overrides := cfg.ApplyTo(&dst)
	assert.Equal(t, "from-cli", dst.Name)
	assert.Equal(t, int32(4), dst.CPUCount)
	assert.True(t, dst.NoCache)
	assert.ElementsMatch(t, []string{"name", "cpu_count"}, overrides)
}

func TestApplyTo_MixedFill(t *testing.T) {
	cfg := &FileConfig{
		Dockerfile: "./Dockerfile",
		StartCmd:   "/start.sh",
	}
	dst := BuildFields{
		Name: "from-cli",
	}
	overrides := cfg.ApplyTo(&dst)
	assert.Equal(t, "from-cli", dst.Name)
	assert.Equal(t, "./Dockerfile", dst.Dockerfile)
	assert.Equal(t, "/start.sh", dst.StartCmd)
	assert.Empty(t, overrides)
}

func TestApplyTo_IdenticalValuesNoOverrideReport(t *testing.T) {
	cfg := &FileConfig{Name: "same", CPUCount: 2}
	dst := BuildFields{Name: "same", CPUCount: 2}
	overrides := cfg.ApplyTo(&dst)
	assert.Empty(t, overrides)
}

func TestApplyTo_BoolDefinedFillsZero(t *testing.T) {
	// 文件显式定义 no_cache = true，dst 为零值 → 应用文件值
	cfg := &FileConfig{NoCache: true, defined: map[string]bool{"no_cache": true}}
	dst := BuildFields{NoCache: false, NoCacheChanged: false}
	overrides := cfg.ApplyTo(&dst)
	assert.True(t, dst.NoCache)
	assert.Empty(t, overrides)
}

func TestApplyTo_BoolNotDefinedIgnored(t *testing.T) {
	// 文件未定义 no_cache（未经 Load 构造），即便结构体零值为 false 也不应覆盖
	cfg := &FileConfig{NoCache: false}
	dst := BuildFields{NoCache: true}
	overrides := cfg.ApplyTo(&dst)
	assert.True(t, dst.NoCache, "CLI-true 在文件未定义 no_cache 时保留")
	assert.Empty(t, overrides)
}

func TestApplyTo_BoolDefinedCLITrueWins(t *testing.T) {
	// 文件 no_cache = false，CLI no_cache = true → CLI 胜出并报告覆盖
	cfg := &FileConfig{NoCache: false, defined: map[string]bool{"no_cache": true}}
	dst := BuildFields{NoCache: true, NoCacheChanged: true}
	overrides := cfg.ApplyTo(&dst)
	assert.True(t, dst.NoCache)
	assert.ElementsMatch(t, []string{"no_cache"}, overrides)
}

func TestApplyTo_BoolDefinedCLIFalseWins(t *testing.T) {
	// 文件 no_cache = true，CLI 显式 no_cache=false → CLI 胜出并报告覆盖。
	cfg := &FileConfig{NoCache: true, defined: map[string]bool{"no_cache": true}}
	dst := BuildFields{NoCache: false, NoCacheChanged: true}
	overrides := cfg.ApplyTo(&dst)
	assert.False(t, dst.NoCache)
	assert.ElementsMatch(t, []string{"no_cache"}, overrides)
}

func TestWriteTemplateID_ReplaceExisting(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := `# comment kept
template_id = ""
name = "demo"
dockerfile = "./Dockerfile"
`
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	expected := `# comment kept
template_id = "tmpl-new"
name = "demo"
dockerfile = "./Dockerfile"
`
	assert.Equal(t, expected, string(got))
}

func TestWriteTemplateID_InsertWhenMissing(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := `name = "demo"
dockerfile = "./Dockerfile"
`
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Contains(t, string(got), `template_id = "tmpl-new"`)
	assert.Contains(t, string(got), `name = "demo"`)
}

func TestWriteTemplateID_PreservesCommentedLine(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := `# template_id = "old"
name = "demo"
`
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Contains(t, string(got), `# template_id = "old"`)
	assert.Contains(t, string(got), `template_id = "tmpl-new"`)
}

func TestWriteTemplateID_PreservesIndentation(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := "    template_id = \"old\"\nname = \"demo\"\n"
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Contains(t, string(got), "    template_id = \"tmpl-new\"")
}

func TestWriteTemplateID_DoesNotReplaceTableField(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := "[section]\n    template_id = \"section-value\"\nname = \"demo\"\n"
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, "template_id = \"tmpl-new\"\n[section]\n    template_id = \"section-value\"\nname = \"demo\"\n", string(got))
}

func TestWriteTemplateID_PreservesCRLF(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	orig := "template_id = \"\"\r\nname = \"demo\"\r\n"
	require.NoError(t, os.WriteFile(p, []byte(orig), 0644))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, "template_id = \"tmpl-new\"\r\nname = \"demo\"\r\n", string(got))
}

func TestWriteTemplateID_PreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(p, []byte("name = \"demo\"\n"), 0600))

	err := WriteTemplateID(p, "tmpl-new")
	require.NoError(t, err)

	info, err := os.Stat(p)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
