package operations

import (
	"context"
	"errors"
	"net/http"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

// lookupTemplateIDByName 在当前环境中按 name 查找并返回 template_id。
// 后端保证同一环境内 alias 唯一，因此可作为稳定定位 key 使用。
//
// 返回值约定：
//   - 找到：返回 template_id, nil
//   - 未找到：返回 "", nil（由调用方决定回退到 create）
//   - 调用 GetTemplateByAlias 出错：返回 "", err
func lookupTemplateIDByName(ctx context.Context, client *sandbox.Client, name string) (string, error) {
	if name == "" {
		return "", nil
	}
	tmpl, err := client.GetTemplateByAlias(ctx, name)
	if err != nil {
		if isTemplateAliasNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return tmpl.TemplateID, nil
}

func isTemplateAliasNotFound(err error) bool {
	var apiErr *sandbox.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
