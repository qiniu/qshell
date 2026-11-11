package operations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
	templatedockerfile "github.com/qiniu/qshell/v2/iqshell/sandbox/template/dockerfile"
)

func TestOwnershipError_RejectsForeignPublicTemplate(t *testing.T) {
	err := ownershipError("agents-base", "cbz63h4i1xiux811uoux", false)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "belongs to another team")
		assert.Contains(t, err.Error(), "unique template name")
	}
}

func TestOwnershipError_AllowsOwnedTemplate(t *testing.T) {
	assert.NoError(t, ownershipError("demo", "tmpl-x", true))
}

func TestLegacyDiskSizeWarning_PointsToMinFreeDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prod.toml")
	require.NoError(t, os.WriteFile(path, []byte("disk_size_mb = 20480\n"), 0o644))
	cfg, err := config.Load(path)
	require.NoError(t, err)

	warning := legacyDiskSizeWarning(cfg)

	assert.Contains(t, warning, "disk_size_mb = 20480")
	assert.Contains(t, warning, path)
	assert.Contains(t, warning, "min_free_disk_mb")
}

func TestLegacyDiskSizeWarning_QuietWithoutLegacyField(t *testing.T) {
	assert.Empty(t, legacyDiskSizeWarning(nil))
	assert.Empty(t, legacyDiskSizeWarning(&config.FileConfig{}))
}

func TestBuildParamsFromDockerfileResult_UsesFromTemplateAsBase(t *testing.T) {
	result := &templatedockerfile.ConvertResult{
		BaseImage: "ignored-from-dockerfile",
		Steps: []sandbox.TemplateStep{
			{Type: "RUN", Args: stringSlicePtr("echo hello")},
		},
	}

	params := buildParamsFromDockerfileResult(result, BuildInfo{
		FromTemplate: "agents-base",
	})

	assert.Nil(t, params.FromImage)
	assert.NotNil(t, params.FromTemplate)
	assert.Equal(t, "agents-base", *params.FromTemplate)
	assert.NotNil(t, params.Steps)
	assert.Len(t, *params.Steps, 1)
}

func TestBuildParamsFromDockerfileResult_FromImageOverridesDockerfileBase(t *testing.T) {
	result := &templatedockerfile.ConvertResult{
		BaseImage: "dockerfile-base",
		Steps:     []sandbox.TemplateStep{},
	}

	params := buildParamsFromDockerfileResult(result, BuildInfo{
		FromImage: "explicit-base",
	})

	assert.NotNil(t, params.FromImage)
	assert.Equal(t, "explicit-base", *params.FromImage)
	assert.Nil(t, params.FromTemplate)
}

func TestCreateTemplateParamsIncludesMinFreeDisk(t *testing.T) {
	minFreeDiskMB := int32(20480)
	params := createTemplateParams(BuildInfo{
		Name:          "demo",
		CPUCount:      2,
		MemoryMB:      2048,
		MinFreeDiskMB: &minFreeDiskMB,
	})

	assert.NotNil(t, params.Name)
	assert.Equal(t, "demo", *params.Name)
	assert.NotNil(t, params.CPUCount)
	assert.NotNil(t, params.MemoryMB)
	assert.NotNil(t, params.MinFreeDiskMB)
	assert.Equal(t, int32(20480), *params.MinFreeDiskMB)
}

func TestCreateTemplateParams_ForwardsCPUAndMemory(t *testing.T) {
	params := createTemplateParams(BuildInfo{
		Name:     "demo",
		CPUCount: 4,
		MemoryMB: 8192,
	})

	require.NotNil(t, params.CPUCount)
	require.NotNil(t, params.MemoryMB)
	assert.Equal(t, int32(4), *params.CPUCount)
	assert.Equal(t, int32(8192), *params.MemoryMB)
}

func TestCreateTemplateParams_OmitsMinFreeDiskWhenUnset(t *testing.T) {
	params := createTemplateParams(BuildInfo{Name: "demo"})

	assert.Nil(t, params.MinFreeDiskMB)
}

func TestCreateTemplateParams_SendsExplicitZero(t *testing.T) {
	zero := int32(0)

	params := createTemplateParams(BuildInfo{Name: "demo", MinFreeDiskMB: &zero})

	require.NotNil(t, params.MinFreeDiskMB)
	assert.Equal(t, int32(0), *params.MinFreeDiskMB)
}

func TestBuildParamsFromDockerfileResult_FallsBackToDockerfileBase(t *testing.T) {
	result := &templatedockerfile.ConvertResult{
		BaseImage: "ubuntu:22.04",
		Steps:     []sandbox.TemplateStep{},
	}

	params := buildParamsFromDockerfileResult(result, BuildInfo{})

	require.NotNil(t, params.FromImage)
	assert.Equal(t, "ubuntu:22.04", *params.FromImage)
	assert.Nil(t, params.FromTemplate)
}

func TestBuildParamsFromDockerfileResult_KeepsSteps(t *testing.T) {
	result := &templatedockerfile.ConvertResult{
		BaseImage: "ubuntu:22.04",
		Steps: []sandbox.TemplateStep{
			{Type: "RUN", Args: stringSlicePtr("echo one")},
			{Type: "RUN", Args: stringSlicePtr("echo two")},
		},
	}

	params := buildParamsFromDockerfileResult(result, BuildInfo{FromImage: "ubuntu:22.04"})

	require.NotNil(t, params.Steps)
	assert.Len(t, *params.Steps, 2)
}

func TestResolveBaseTemplate_EmptyWhenNoBase(t *testing.T) {
	base, err := resolveBaseTemplate("", "")

	require.NoError(t, err)
	assert.Empty(t, base)
}

func TestValidateBuildSourceSelection_AllowsDockerfileOnly(t *testing.T) {
	assert.NoError(t, validateBuildSourceSelection(BuildInfo{Dockerfile: "./Dockerfile"}))
}

func TestFormatBuildFailureLogs(t *testing.T) {
	got := formatBuildFailureLogs([]string{"step 1 failed", "command exited with status 1"})
	assert.Equal(t, "\nBuild Logs:\n  step 1 failed\n  command exited with status 1\n", got)
}

func TestFormatBuildFailureLogs_Empty(t *testing.T) {
	assert.Empty(t, formatBuildFailureLogs(nil))
}

func TestValidateBuildSourceSelection_RejectsFromImageAndFromTemplate(t *testing.T) {
	err := validateBuildSourceSelection(BuildInfo{
		FromImage:    "ubuntu:22.04",
		FromTemplate: "agents-base",
	})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "cannot use both")
	}
}

func TestValidateBuildSourceSelection_AllowsDockerfileWithFromTemplate(t *testing.T) {
	err := validateBuildSourceSelection(BuildInfo{
		Dockerfile:   "./Dockerfile",
		FromTemplate: "agents-base",
	})

	assert.NoError(t, err)
}

func TestResolveBaseTemplate_UsesTemplateIDAsBase(t *testing.T) {
	base, err := resolveBaseTemplate("tmpl-xxxxxxxxxxxx", "")

	require.NoError(t, err)
	assert.Equal(t, "tmpl-xxxxxxxxxxxx", base)
}

func TestResolveBaseTemplate_UsesFromTemplateAsBase(t *testing.T) {
	base, err := resolveBaseTemplate("", "agents-base")

	require.NoError(t, err)
	assert.Equal(t, "agents-base", base)
}

func TestResolveBaseTemplate_RejectsBothSources(t *testing.T) {
	_, err := resolveBaseTemplate("tmpl-xxxxxxxxxxxx", "agents-base")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--template-id")
		assert.Contains(t, err.Error(), "--from-template")
	}
}

func TestCLIBaseSelection_RejectsMultipleFlags(t *testing.T) {
	_, _, err := cliBaseSelection(BuildInfo{
		TemplateID:   "tmpl-base",
		FromTemplate: "agents-base",
	})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--template-id")
		assert.Contains(t, err.Error(), "--from-template")
	}
}

func TestCLIBaseSelection_RejectsImageAndTemplate(t *testing.T) {
	_, _, err := cliBaseSelection(BuildInfo{FromImage: "ubuntu:22.04", FromTemplate: "agents-base"})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--from-image")
		assert.Contains(t, err.Error(), "--from-template")
	}
}

func TestCLIBaseSelection_ReturnsTemplateID(t *testing.T) {
	field, value, err := cliBaseSelection(BuildInfo{TemplateID: "tmpl-base"})

	require.NoError(t, err)
	assert.Equal(t, "template_id", field)
	assert.Equal(t, "tmpl-base", value)
}

func TestCLIBaseSelection_FallsBackToFromImage(t *testing.T) {
	field, value, err := cliBaseSelection(BuildInfo{FromImage: "ubuntu:22.04"})

	require.NoError(t, err)
	assert.Equal(t, "from_image", field)
	assert.Equal(t, "ubuntu:22.04", value)
}

func TestCLIBaseSelection_EmptyWithoutBase(t *testing.T) {
	field, value, err := cliBaseSelection(BuildInfo{Name: "demo"})

	require.NoError(t, err)
	assert.Empty(t, field)
	assert.Empty(t, value)
}

func TestApplyCLIBasePrecedence_DropsFileBase(t *testing.T) {
	info := BuildInfo{TemplateID: "tmpl-cli", FromTemplate: "file-base"}

	overrode := applyCLIBasePrecedence(&info, "template_id", "tmpl-cli")

	assert.True(t, overrode)
	assert.Equal(t, "tmpl-cli", info.TemplateID)
	assert.Empty(t, info.FromTemplate)
}

func TestApplyCLIBasePrecedence_KeepsCLIOnly(t *testing.T) {
	info := BuildInfo{FromTemplate: "cli-base"}

	overrode := applyCLIBasePrecedence(&info, "from_template", "cli-base")

	assert.False(t, overrode)
	assert.Equal(t, "cli-base", info.FromTemplate)
}

func TestRejectSelfBase_BlamesFlagWhenBaseCameFromCLI(t *testing.T) {
	err := rejectSelfBase(BuildInfo{Name: "demo", FromTemplate: "tmpl-demo"}, "tmpl-demo", nil, "template_id")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--template-id")
		assert.NotContains(t, err.Error(), "qshell.sandbox.toml")
	}
}

func TestBaseConflictError_ReportsConfigKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(path, []byte("name = \"demo\"\ntemplate_id = \"tmpl-old\"\nfrom_template = \"base\"\n"), 0o644))
	cfg, err := config.Load(path)
	require.NoError(t, err)

	conflict := baseConflictError(cfg, "")

	if assert.Error(t, conflict) {
		assert.Contains(t, conflict.Error(), path)
		assert.Contains(t, conflict.Error(), "template_id")
		assert.Contains(t, conflict.Error(), "from_template")
	}
}

func TestBaseConflictError_IgnoresFileWhenCLIWins(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(path, []byte("template_id = \"a\"\nfrom_template = \"b\"\n"), 0o644))
	cfg, err := config.Load(path)
	require.NoError(t, err)

	assert.NoError(t, baseConflictError(cfg, "template_id"))
}

func TestRejectSelfBase_RejectsTemplateIDItself(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(path, []byte("name = \"demo\"\ntemplate_id = \"tmpl-demo\"\n"), 0o644))
	cfg, err := config.Load(path)
	require.NoError(t, err)

	buildErr := rejectSelfBase(BuildInfo{Name: "demo", FromTemplate: "tmpl-demo"}, "tmpl-demo", cfg, "")

	if assert.Error(t, buildErr) {
		assert.Contains(t, buildErr.Error(), path)
		assert.Contains(t, buildErr.Error(), "remove that line")
	}
}

func TestRejectSelfBase_AllowsDistinctBase(t *testing.T) {
	assert.NoError(t, rejectSelfBase(BuildInfo{Name: "demo", FromTemplate: "base"}, "tmpl-demo", nil, ""))
}

func loadTestConfig(t *testing.T, content string) *config.FileConfig {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	cfg, err := config.Load(path)
	require.NoError(t, err)
	return cfg
}

func TestResolveBuildInputs_ReportsTOMLBaseConflict(t *testing.T) {
	cfg := loadTestConfig(t, "name = \"demo\"\ntemplate_id = \"tmpl-old\"\nfrom_template = \"base\"\n")

	_, _, err := resolveBuildInputs(BuildInfo{}, cfg, "", "")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), cfg.SourcePath())
		assert.Contains(t, err.Error(), "template_id")
		assert.Contains(t, err.Error(), "from_template")
	}
}

func TestResolveBuildInputs_UsesFileFromTemplate(t *testing.T) {
	cfg := loadTestConfig(t, "name = \"demo\"\nfrom_template = \"base\"\n")

	info, overrode, err := resolveBuildInputs(BuildInfo{Name: "demo", FromTemplate: "base"}, cfg, "", "")

	require.NoError(t, err)
	assert.False(t, overrode)
	assert.Equal(t, "base", info.FromTemplate)
	assert.Empty(t, info.TemplateID)
}

func TestResolveBuildInputs_CLIBaseOverridesFile(t *testing.T) {
	cfg := loadTestConfig(t, "name = \"demo\"\nfrom_template = \"file-base\"\n")

	info, overrode, err := resolveBuildInputs(
		BuildInfo{Name: "demo", TemplateID: "tmpl-cli", FromTemplate: "file-base"}, cfg, "template_id", "tmpl-cli")

	require.NoError(t, err)
	assert.True(t, overrode)
	assert.Equal(t, "tmpl-cli", info.FromTemplate)
}

func TestResolveBuildInputs_LegacyTemplateIDNeedsName(t *testing.T) {
	cfg := loadTestConfig(t, "template_id = \"tmpl-old\"\n")

	_, _, err := resolveBuildInputs(BuildInfo{TemplateID: "tmpl-old"}, cfg, "", "")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "add a name")
		assert.Contains(t, err.Error(), "template_id")
	}
}

func TestResolveBuildInputs_RequiresName(t *testing.T) {
	_, _, err := resolveBuildInputs(BuildInfo{FromTemplate: "base"}, nil, "", "")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "required")
	}
}

func TestResolveBuildInputs_RejectsSelfBase(t *testing.T) {
	_, _, err := resolveBuildInputs(BuildInfo{Name: "demo", FromTemplate: "demo"}, nil, "", "")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "own base")
	}
}

func TestResolveBuildInputs_CLIFromImageOverridesFileTemplate(t *testing.T) {
	cfg := loadTestConfig(t, "name = \"demo\"\nfrom_template = \"file-base\"\n")

	info, overrode, err := resolveBuildInputs(
		BuildInfo{Name: "demo", FromImage: "ubuntu:22.04", FromTemplate: "file-base"}, cfg, "from_image", "ubuntu:22.04")

	require.NoError(t, err)
	assert.True(t, overrode)
	assert.Equal(t, "ubuntu:22.04", info.FromImage)
	assert.Empty(t, info.FromTemplate)
}

func TestCliBaseSelection_NamingMatchesSelfBaseGuard(t *testing.T) {
	field, value, err := cliBaseSelection(BuildInfo{TemplateID: "tmpl-demo"})
	require.NoError(t, err)

	guardErr := rejectSelfBase(BuildInfo{Name: "demo", FromTemplate: value}, value, nil, field)

	if assert.Error(t, guardErr) {
		assert.Contains(t, guardErr.Error(), "--template-id names the template being built")
	}
}

func TestPrepareDockerfileBuild_ParsesSteps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Dockerfile")
	require.NoError(t, os.WriteFile(path, []byte("FROM ubuntu:22.04\nRUN echo hi\n"), 0o644))

	plan, err := prepareDockerfileBuild(BuildInfo{Dockerfile: path})

	require.NoError(t, err)
	require.NotNil(t, plan)
	assert.NotEmpty(t, plan.result.Steps)
}

func TestPrepareDockerfileBuild_MissingFile(t *testing.T) {
	_, err := prepareDockerfileBuild(BuildInfo{Dockerfile: filepath.Join(t.TempDir(), "missing")})

	assert.Error(t, err)
}

func TestValidateBuildSourceSelection_RequiresSource(t *testing.T) {
	err := validateBuildSourceSelection(BuildInfo{})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--from-image")
	}
}

func stringSlicePtr(values ...string) *[]string {
	return &values
}

func TestMinDiskConfigToCreateParams(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		cli     *int32
		want    *int32
	}{
		{name: "unset"},
		{name: "configured", content: "min_disk_mb = 10240", want: diskValue(10240)},
		{name: "configured zero", content: "min_disk_mb = 0", want: diskValue(0)},
		{name: "CLI zero overrides", content: "min_disk_mb = 10240", cli: diskValue(0), want: diskValue(0)},
		{name: "CLI positive overrides", content: "min_disk_mb = 10240", cli: diskValue(20480), want: diskValue(20480)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "qshell.sandbox.toml")
			require.NoError(t, os.WriteFile(path, []byte(tc.content), 0600))
			info, _, err := loadBuildConfig(BuildInfo{ConfigPath: path, MinDiskMB: tc.cli, MinFreeDiskMB: diskValue(5120)})
			require.NoError(t, err)
			params := createTemplateParams(info)
			assert.Equal(t, tc.want, params.MinDiskMB)
			assert.Equal(t, diskValue(5120), params.MinFreeDiskMB)
		})
	}
}

func diskValue(value int32) *int32 { return &value }
