package operations

import (
	"fmt"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
)

// templateIDFromCwdConfig 在当前工作目录加载 qshell.sandbox.toml，
// 返回其中的 template_id。加载/解析失败时打印错误并返回 (""，false)。
// 文件不存在或未包含 template_id 时返回 ("", true)，由调用方处理。
func templateIDFromCwdConfig() (string, bool) {
	cfg, err := config.LoadFromCwd()
	if err != nil {
		sbClient.PrintError("load config: %v", err)
		return "", false
	}
	if cfg == nil || cfg.TemplateID == "" {
		return "", true
	}
	fmt.Fprintf(os.Stderr, "[config] using template_id from %s\n", cfg.SourcePath())
	return cfg.TemplateID, true
}
