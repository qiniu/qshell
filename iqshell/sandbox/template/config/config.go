// Package config 处理 qshell.sandbox.toml 配置文件的加载、合并与回写。
package config

// DefaultFileName 是模板配置文件的默认文件名。
const DefaultFileName = "qshell.sandbox.toml"

// FileConfig 对应 qshell.sandbox.toml 中的所有可配置字段。
// 所有字段均为可选：未设置时回退到 CLI 参数或内置默认值。
type FileConfig struct {
	// TemplateID 是模板 ID，首次构建后由 qshell 自动回写。
	TemplateID string `toml:"template_id"`

	// Name 是模板名称，仅首次构建（未设置 TemplateID 时）生效。
	Name string `toml:"name"`

	// Dockerfile 是 Dockerfile 路径，启用 v2 构建。
	Dockerfile string `toml:"dockerfile"`

	// Path 是构建上下文目录，默认为 Dockerfile 所在目录。
	Path string `toml:"path"`

	// FromImage 是基础 Docker 镜像。
	FromImage string `toml:"from_image"`

	// FromTemplate 是基础模板。
	FromTemplate string `toml:"from_template"`

	// StartCmd 是容器启动命令。
	StartCmd string `toml:"start_cmd"`

	// ReadyCmd 是就绪检查命令。
	ReadyCmd string `toml:"ready_cmd"`

	// CPUCount 是沙箱 CPU 核数。
	CPUCount int32 `toml:"cpu_count"`

	// MemoryMB 是沙箱内存大小（MiB）。
	MemoryMB int32 `toml:"memory_mb"`

	// NoCache 强制完整构建，忽略缓存。
	NoCache bool `toml:"no_cache"`

	// sourcePath 记录配置文件的绝对路径，用于回写。未导出。
	sourcePath string

	// defined 记录 TOML 文件中显式定义的字段 key 集合。未导出。
	defined map[string]bool
}

// SourcePath 返回配置来源文件的绝对路径，未从文件加载时返回空字符串。
func (c *FileConfig) SourcePath() string {
	return c.sourcePath
}
