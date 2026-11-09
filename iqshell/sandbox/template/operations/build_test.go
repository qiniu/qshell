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

func TestShouldRecreateForDisk_RecreatesResolvedNameWithFreeDisk(t *testing.T) {
	minFree := int32(20480)

	assert.True(t, shouldRecreateForDisk(BuildInfo{Name: "agents-base", MinFreeDiskMB: &minFree}, true))
}

func TestShouldRecreateForDisk_KeepsRebuildWithoutDiskSetting(t *testing.T) {
	minFree := int32(20480)

	assert.False(t, shouldRecreateForDisk(BuildInfo{Name: "agents-base"}, true))
	assert.False(t, shouldRecreateForDisk(BuildInfo{Name: "agents-base", MinFreeDiskMB: &minFree}, false))
	assert.False(t, shouldRecreateForDisk(BuildInfo{MinFreeDiskMB: &minFree}, true))
}

func TestShouldRecreateForDisk_AppliesExplicitZero(t *testing.T) {
	zero := int32(0)

	assert.True(t, shouldRecreateForDisk(BuildInfo{Name: "agents-base", MinFreeDiskMB: &zero}, true))
}

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
		assert.Contains(t, err.Error(), "cannot specify both")
	}
}

func TestValidateBuildSourceSelection_AllowsDockerfileWithFromTemplate(t *testing.T) {
	err := validateBuildSourceSelection(BuildInfo{
		Dockerfile:   "./Dockerfile",
		FromTemplate: "agents-base",
	})

	assert.NoError(t, err)
}

func TestValidateRebuildSourceSelection_RejectsCLIFromImage(t *testing.T) {
	info := BuildInfo{
		TemplateID: "tmpl-xxxxxxxxxxxx",
		Dockerfile: "./Dockerfile",
		FromImage:  "ubuntu:22.04",
	}
	err := validateRebuildSourceSelection(info, true, false)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--from-image")
		assert.Contains(t, err.Error(), "--from-template")
	}
}

func TestValidateRebuildSourceSelection_RejectsCLIFromTemplate(t *testing.T) {
	info := BuildInfo{
		TemplateID:   "tmpl-xxxxxxxxxxxx",
		Dockerfile:   "./Dockerfile",
		FromTemplate: "agents-base",
	}
	err := validateRebuildSourceSelection(info, false, true)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "--from-image")
		assert.Contains(t, err.Error(), "--from-template")
	}
}

// TOML 配置中声明的 from_image 在 rebuild 时必须保留，
// 否则下游 buildFromDockerfile 会回退到 Dockerfile 的 FROM 行。
func TestValidateRebuildSourceSelection_PreservesConfigFromImage(t *testing.T) {
	info := BuildInfo{
		TemplateID: "tmpl-xxxxxxxxxxxx",
		Dockerfile: "./Dockerfile",
		FromImage:  "ubuntu:22.04",
	}

	err := validateRebuildSourceSelection(info, false, false)

	assert.NoError(t, err)
	assert.Equal(t, "ubuntu:22.04", info.FromImage)
	assert.Empty(t, info.FromTemplate)
}

// TOML 配置中声明的 from_template 在 rebuild 时必须保留，
// 这是修复 "image 'scratch' not found" 的关键场景：
// agents-base 等模板基于 from_template = "base" 构建，
// rebuild 时若被清空会让 Dockerfile 中的 FROM scratch 直接送到 builder。
func TestValidateRebuildSourceSelection_PreservesConfigFromTemplate(t *testing.T) {
	info := BuildInfo{
		TemplateID:   "tmpl-xxxxxxxxxxxx",
		Dockerfile:   "./Dockerfile",
		FromTemplate: "base",
	}

	err := validateRebuildSourceSelection(info, false, false)

	assert.NoError(t, err)
	assert.Empty(t, info.FromImage)
	assert.Equal(t, "base", info.FromTemplate)
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
