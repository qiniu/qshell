package config

// BuildFields 是 operations.BuildInfo 中可由配置文件提供的字段子集。
// 单独定义以避免 config 包依赖 operations 包（解决循环依赖）。
// operations.Build 在调用前将 BuildInfo 拷贝为 BuildFields，调用后拷贝回。
type BuildFields struct {
	// TemplateID 是模板 ID。
	TemplateID string
	// Name 是模板名称。
	Name string
	// Dockerfile 是 Dockerfile 路径。
	Dockerfile string
	// Path 是构建上下文目录。
	Path string
	// FromImage 是基础 Docker 镜像。
	FromImage string
	// FromTemplate 是基础模板。
	FromTemplate string
	// StartCmd 是容器启动命令。
	StartCmd string
	// ReadyCmd 是就绪检查命令。
	ReadyCmd string
	// CPUCount 是沙箱 CPU 核数。
	CPUCount int32
	// MemoryMB 是沙箱内存大小（MiB）。
	MemoryMB int32
	// NoCache 强制完整构建，忽略缓存。
	NoCache bool
	// NoCacheChanged 表示 CLI 是否显式设置了 --no-cache。
	NoCacheChanged bool
}

// ApplyTo 将 FileConfig 的值合并到 dst 中：仅当 dst 字段为零值时才采用 file 值。
// 返回被 CLI 覆盖（即 dst 已有非零值且与 file 值不同）的 TOML 字段名列表。
func (c *FileConfig) ApplyTo(dst *BuildFields) []string {
	var overrides []string

	applyString := func(fileVal string, dstVal *string, key string) {
		if fileVal == "" {
			return
		}
		if *dstVal == "" {
			*dstVal = fileVal
			return
		}
		if *dstVal != fileVal {
			overrides = append(overrides, key)
		}
	}
	applyInt32 := func(fileVal int32, dstVal *int32, key string) {
		if fileVal == 0 {
			return
		}
		if *dstVal == 0 {
			*dstVal = fileVal
			return
		}
		if *dstVal != fileVal {
			overrides = append(overrides, key)
		}
	}
	applyBool := func(fileVal bool, dstVal *bool, key string) {
		// 未在文件中显式定义时，尊重 CLI 值。
		if !c.defined[key] {
			return
		}
		if *dstVal == fileVal {
			return
		}
		if dst.NoCacheChanged {
			overrides = append(overrides, key)
			return
		}
		if !*dstVal {
			*dstVal = fileVal
			return
		}
		overrides = append(overrides, key)
	}

	applyString(c.TemplateID, &dst.TemplateID, "template_id")
	applyString(c.Name, &dst.Name, "name")
	applyString(c.Dockerfile, &dst.Dockerfile, "dockerfile")
	applyString(c.Path, &dst.Path, "path")
	applyString(c.FromImage, &dst.FromImage, "from_image")
	applyString(c.FromTemplate, &dst.FromTemplate, "from_template")
	applyString(c.StartCmd, &dst.StartCmd, "start_cmd")
	applyString(c.ReadyCmd, &dst.ReadyCmd, "ready_cmd")
	applyInt32(c.CPUCount, &dst.CPUCount, "cpu_count")
	applyInt32(c.MemoryMB, &dst.MemoryMB, "memory_mb")
	applyBool(c.NoCache, &dst.NoCache, "no_cache")

	return overrides
}
