package operations

import (
	"context"
	"fmt"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
)

// templateIDFromCwdConfig 在当前工作目录加载 qshell.sandbox.toml，
// 返回其中可用于定位模板的 template_id。
//
// 解析顺序：
//  1. toml 中的 template_id（向后兼容已有配置）
//  2. toml 中的 name → 调用 ListTemplates 按 alias 查找
//
// 这样 toml 里只写 name（不写 template_id）也能让 publish/get/delete
// 等命令工作，便于多环境共享同一份配置。
//
// 加载/解析失败或远端调用失败时打印错误并返回 (""，false)。
// 配置文件不存在或既无 template_id 也无 name 时返回 ("", true)，
// 由调用方决定后续行为。
func templateIDFromCwdConfig() (string, bool) {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		sbClient.PrintError("load config: %v", err)
		return "", false
	}
	if cfg == nil {
		return "", true
	}
	if cfg.TemplateID != "" {
		fmt.Fprintf(os.Stderr, "[config] using template_id from %s\n", cfg.SourcePath())
		return cfg.TemplateID, true
	}
	if cfg.Name == "" {
		return "", true
	}

	client, cErr := sbClient.NewSandboxClient()
	if cErr != nil {
		sbClient.PrintError("%v", cErr)
		return "", false
	}
	id, lErr := lookupTemplateIDByName(context.Background(), client, cfg.Name)
	if lErr != nil {
		sbClient.PrintError("lookup template by name %q: %v", cfg.Name, lErr)
		return "", false
	}
	if id == "" {
		return "", true
	}
	fmt.Fprintf(os.Stderr, "[lookup] template %q resolved to %s (from %s)\n",
		cfg.Name, id, cfg.SourcePath())
	return id, true
}
