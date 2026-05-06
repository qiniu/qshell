package operations

import (
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
)

func TestFindTemplateIDByAlias_Hit(t *testing.T) {
	templates := []sandbox.Template{
		{TemplateID: "id-base", Aliases: []string{"agents-base"}},
		{TemplateID: "id-claude", Aliases: []string{"claude"}},
		{TemplateID: "id-codex", Aliases: []string{"codex"}},
	}

	assert.Equal(t, "id-claude", findTemplateIDByAlias(templates, "claude"))
	assert.Equal(t, "id-base", findTemplateIDByAlias(templates, "agents-base"))
}

func TestFindTemplateIDByAlias_MultipleAliases(t *testing.T) {
	templates := []sandbox.Template{
		{TemplateID: "id-claude", Aliases: []string{"claude", "claude-stable"}},
	}

	assert.Equal(t, "id-claude", findTemplateIDByAlias(templates, "claude-stable"))
}

func TestFindTemplateIDByAlias_Miss(t *testing.T) {
	templates := []sandbox.Template{
		{TemplateID: "id-base", Aliases: []string{"agents-base"}},
	}

	assert.Equal(t, "", findTemplateIDByAlias(templates, "claude"))
}

func TestFindTemplateIDByAlias_EmptyName(t *testing.T) {
	templates := []sandbox.Template{
		{TemplateID: "id-base", Aliases: []string{"agents-base"}},
	}

	// 防御性：空 name 不应匹配任何模板，避免误返回第一个 alias 为空字符串的模板。
	assert.Equal(t, "", findTemplateIDByAlias(templates, ""))
}

func TestFindTemplateIDByAlias_EmptyTemplates(t *testing.T) {
	assert.Equal(t, "", findTemplateIDByAlias(nil, "claude"))
	assert.Equal(t, "", findTemplateIDByAlias([]sandbox.Template{}, "claude"))
}
