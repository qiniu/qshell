package operations

import (
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"

	templatedockerfile "github.com/qiniu/qshell/v2/iqshell/sandbox/template/dockerfile"
)

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
