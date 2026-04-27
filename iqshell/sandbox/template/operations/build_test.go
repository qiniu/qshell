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

func TestNormalizeRebuildSourceSelection_RejectsCLIFromImage(t *testing.T) {
	info := BuildInfo{
		TemplateID: "tmpl-xxxxxxxxxxxx",
		Dockerfile: "./Dockerfile",
		FromImage:  "ubuntu:22.04",
	}
	err := normalizeRebuildSourceSelection(&info, true, false)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "when rebuilding")
	}
}

func TestNormalizeRebuildSourceSelection_RejectsCLIFromTemplate(t *testing.T) {
	info := BuildInfo{
		TemplateID:   "tmpl-xxxxxxxxxxxx",
		Dockerfile:   "./Dockerfile",
		FromTemplate: "agents-base",
	}
	err := normalizeRebuildSourceSelection(&info, false, true)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "when rebuilding")
	}
}

func TestNormalizeRebuildSourceSelection_IgnoresConfigFromImage(t *testing.T) {
	info := BuildInfo{
		TemplateID: "tmpl-xxxxxxxxxxxx",
		Dockerfile: "./Dockerfile",
		FromImage:  "ubuntu:22.04",
	}

	err := normalizeRebuildSourceSelection(&info, false, false)

	assert.NoError(t, err)
	assert.Empty(t, info.FromImage)
	assert.Empty(t, info.FromTemplate)
}

func TestNormalizeRebuildSourceSelection_IgnoresConfigFromTemplate(t *testing.T) {
	info := BuildInfo{
		TemplateID:   "tmpl-xxxxxxxxxxxx",
		Dockerfile:   "./Dockerfile",
		FromTemplate: "agents-base",
	}

	err := normalizeRebuildSourceSelection(&info, false, false)

	assert.NoError(t, err)
	assert.Empty(t, info.FromImage)
	assert.Empty(t, info.FromTemplate)
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
