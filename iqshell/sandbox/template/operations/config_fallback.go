package operations

import (
	"context"
	"fmt"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
)

// templateIDFromCwdConfig 在当前工作目录加载 qshell.sandbox.toml，
// 按 name 解析出用于定位模板的 template_id。
//
// template_id 在配置中表示构建所用底模板（与 from_template 等价），
// 因此不再作为 get/delete/publish/unpublish 的目标；若配置里只有
// template_id 而没有 name，直接报错并提示补充 name。
//
// 配置文件不存在或既无 name 也无 template_id 时返回 ("", true)，
// 由调用方决定后续行为；加载/解析或远端调用失败时返回 ("", false)。
func templateIDFromCwdConfig() (string, bool) {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		sbClient.PrintError("load config: %v", err)
		return "", false
	}
	if cfg == nil {
		return "", true
	}
	if cfg.Name == "" {
		if cfg.TemplateID != "" {
			sbClient.PrintError("%s: template_id is the base template for build; add name to locate the template",
				cfg.SourcePath())
			return "", false
		}
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
