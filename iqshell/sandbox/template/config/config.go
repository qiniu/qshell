// Package config 处理 qshell.sandbox.toml 配置文件的加载与合并。
package config

// DefaultFileName 是模板配置文件的默认文件名。
const DefaultFileName = "qshell.sandbox.toml"

// FileConfig 对应 qshell.sandbox.toml 中的所有可配置字段。
// 所有字段均为可选：未设置时回退到 CLI 参数或内置默认值。
type FileConfig struct {
	// TemplateID 是构建所用的底模板 ID，与 from_template 等价。
	TemplateID string `toml:"template_id"`

	// Name 是模板名称，每次构建都按它定位目标模板。
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

	// MinDiskMB 是构建后的最小总磁盘大小（MiB），nil 或 0 表示不设置最小值。
	MinDiskMB *int32 `toml:"min_disk_mb"`

	// MinFreeDiskMB 是构建后的最小剩余空间（MiB），nil 使用团队默认值，0 表示不为剩余空间扩容。
	MinFreeDiskMB *int32 `toml:"min_free_disk_mb"`

	// DiskSizeMB 是已废弃的构建磁盘大小字段（MiB）。
	// 这里仅保留解析能力，以便 build 时提示改用 min_disk_mb
	// 表达最小总磁盘大小，而不是静默忽略用户配置。
	DiskSizeMB *int32 `toml:"disk_size_mb"`

	// NoCache 强制完整构建，忽略缓存。
	NoCache bool `toml:"no_cache"`

	// sourcePath 记录配置文件的绝对路径，用于覆盖提示。未导出。
	sourcePath string

	// defined 记录 TOML 文件中显式定义的字段 key 集合。未导出。
	defined map[string]bool
}

// SourcePath 返回配置来源文件的绝对路径，未从文件加载时返回空字符串。
func (c *FileConfig) SourcePath() string {
	return c.sourcePath
}
