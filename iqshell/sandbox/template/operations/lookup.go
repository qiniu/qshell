package operations

import (
	"context"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

// lookupTemplateIDByName 在当前环境的模板列表中按 name 查找并返回 template_id。
// 通过遍历 ListTemplates 结果、匹配 Aliases 实现：后端保证同一环境内 alias 唯一，
// 因此命中即可作为稳定定位 key 使用。
//
// 返回值约定：
//   - 找到：返回 template_id, nil
//   - 未找到：返回 "", nil（由调用方决定回退到 create）
//   - 调用 ListTemplates 出错：返回 "", err
func lookupTemplateIDByName(ctx context.Context, client *sandbox.Client, name string) (string, error) {
	if name == "" {
		return "", nil
	}
	templates, err := client.ListTemplates(ctx, nil)
	if err != nil {
		return "", err
	}
	return findTemplateIDByAlias(templates, name), nil
}

// findTemplateIDByAlias 在已有模板切片中按 alias 精确匹配并返回 template_id。
// 未命中返回 ""。提取为纯函数以便单元测试。
func findTemplateIDByAlias(templates []sandbox.Template, name string) string {
	if name == "" {
		return ""
	}
	for _, t := range templates {
		for _, alias := range t.Aliases {
			if alias == name {
				return t.TemplateID
			}
		}
	}
	return ""
}
