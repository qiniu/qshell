package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/dockerfile"
)

// BuildInfo 保存构建模板的参数。
type BuildInfo struct {
	// Name 是模板名称（用于创建 + 构建）。
	Name string

	// TemplateID 是构建所用的底模板 ID，与 --from-template 等价。
	TemplateID string

	// FromImage 是基础 Docker 镜像。
	FromImage string

	// FromTemplate 是基础模板。
	FromTemplate string

	// StartCmd 是构建完成后执行的命令。
	StartCmd string

	// ReadyCmd 是就绪检查命令。
	ReadyCmd string

	// CPUCount 是沙箱 CPU 核数。
	CPUCount int32

	// MemoryMB 是沙箱内存大小（MiB）。
	MemoryMB int32

	// MinFreeDiskMB 是构建后的最小剩余空间（MiB），每次构建都会提交。
	MinFreeDiskMB *int32

	// Wait 指示是否等待构建完成。
	Wait bool

	// NoCache 强制完整构建，忽略缓存。
	NoCache bool
	// NoCacheChanged 表示 CLI 是否显式设置了 --no-cache。
	NoCacheChanged bool

	// Dockerfile 是 Dockerfile 的路径（启用 v2 Dockerfile 构建）。
	Dockerfile string

	// Path 是构建上下文目录（默认为 Dockerfile 所在目录）。
	Path string

	// ConfigPath 是 qshell.sandbox.toml 的显式路径。
	// 为空时自动在当前工作目录查找。
	ConfigPath string
}

// buildFailed 打印错误并设置命令失败状态，便于 CI 感知构建失败。
func buildFailed(format string, a ...any) {
	data.SetCmdStatusError()
	sbClient.PrintError(format, a...)
}

// Build 分两步创建模板。
//
// 先通过创建接口申请构建（CPU、内存、最小剩余空间等资源参数在这一步
// 确定），再通过构建接口指定底模板/镜像与构建步骤。
// 同名模板已存在时，创建接口会沿用该模板 ID 并应用新的资源设置。
func Build(info BuildInfo) {
	// 记录 CLI 显式指定的底模板，合并后用它覆盖配置文件里的底模板字段。
	cliBaseField, cliBase, baseErr := cliBaseSelection(info)
	if baseErr != nil {
		buildFailed("%v", baseErr)
		return
	}

	// 加载配置文件并合并（CLI > file > default）
	var cfg *config.FileConfig
	if info.ConfigPath != "" {
		loaded, cErr := config.Load(info.ConfigPath)
		if cErr != nil {
			buildFailed("load config: %v", cErr)
			return
		}
		if loaded == nil {
			buildFailed("config file not found: %s", info.ConfigPath)
			return
		}
		cfg = loaded
	} else {
		loaded, cErr := config.LoadFromCwd()
		if cErr != nil {
			buildFailed("load config: %v", cErr)
			return
		}
		cfg = loaded
	}

	if cfg != nil {
		fields := config.BuildFields{
			TemplateID:     info.TemplateID,
			Name:           info.Name,
			Dockerfile:     info.Dockerfile,
			Path:           info.Path,
			FromImage:      info.FromImage,
			FromTemplate:   info.FromTemplate,
			StartCmd:       info.StartCmd,
			ReadyCmd:       info.ReadyCmd,
			CPUCount:       info.CPUCount,
			MemoryMB:       info.MemoryMB,
			MinFreeDiskMB:  info.MinFreeDiskMB,
			NoCache:        info.NoCache,
			NoCacheChanged: info.NoCacheChanged,
		}
		overrides := cfg.ApplyTo(&fields)
		info.TemplateID = fields.TemplateID
		info.Name = fields.Name
		info.Dockerfile = fields.Dockerfile
		info.Path = fields.Path
		info.FromImage = fields.FromImage
		info.FromTemplate = fields.FromTemplate
		info.StartCmd = fields.StartCmd
		info.ReadyCmd = fields.ReadyCmd
		info.CPUCount = fields.CPUCount
		info.MemoryMB = fields.MemoryMB
		info.MinFreeDiskMB = fields.MinFreeDiskMB
		info.NoCache = fields.NoCache

		for _, key := range overrides {
			fmt.Fprintf(os.Stderr, "[config] CLI overrides %s from %s\n", key, cfg.SourcePath())
		}
	}
	if warning := legacyDiskSizeWarning(cfg); warning != "" {
		sbClient.PrintWarn("%s", warning)
	}
	if cliBase != "" {
		if overrode := applyCLIBasePrecedence(&info, cliBaseField, cliBase); overrode && cfg != nil {
			fmt.Fprintf(os.Stderr, "[config] CLI overrides base template from %s\n", cfg.SourcePath())
		}
	}

	if err := baseConflictError(cfg, cliBase); err != nil {
		buildFailed("%v", err)
		return
	}
	base, bErr := resolveBaseTemplate(info.TemplateID, info.FromTemplate)
	if bErr != nil {
		buildFailed("%v", bErr)
		return
	}
	info.FromTemplate = base
	if base != "" && base == info.Name {
		buildFailed("cannot use template %q as its own base template", info.Name)
		return
	}

	if info.Name == "" {
		buildFailed("template name (--name) is required")
		return
	}
	if err := validateBuildSourceSelection(info); err != nil {
		buildFailed("%v", err)
		return
	}
	if info.MinFreeDiskMB != nil && *info.MinFreeDiskMB < 0 {
		data.SetCmdStatusError()
		buildFailed("--min-free-disk / min_free_disk_mb must be non-negative")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		buildFailed("%v", err)
		return
	}

	ctx := context.Background()

	// 同名模板已存在时只做归属校验：创建接口会沿用该模板 ID 并应用
	// 新的资源设置，但不会接管其他团队的同名公共模板。
	existingID, lookupErr := lookupTemplateIDByName(ctx, client, info.Name)
	if lookupErr != nil {
		buildFailed("lookup template by name %q: %v", info.Name, lookupErr)
		return
	}
	if existingID != "" {
		fmt.Fprintf(os.Stderr, "[lookup] template %q resolved to %s; building under this name\n", info.Name, existingID)
		if err := ensureTemplateOwned(ctx, client, info.Name, existingID); err != nil {
			buildFailed("%v", err)
			return
		}
		if err := rejectSelfBase(info, existingID, cfg, cliBase); err != nil {
			buildFailed("%v", err)
			return
		}
	}

	createParams := createTemplateParams(info)

	action := "Creating"
	if existingID != "" {
		action = "Requesting build for existing"
	}
	fmt.Printf("%s template %s...\n", action, info.Name)
	resp, cErr := client.CreateTemplate(ctx, createParams)
	if cErr != nil {
		buildFailed("create template failed: %v", cErr)
		return
	}
	templateID := resp.TemplateID
	buildID := resp.BuildID
	sbClient.PrintSuccess("Template %s ready for build (build ID: %s)", templateID, buildID)

	if info.Dockerfile != "" {
		if err := buildFromDockerfile(ctx, client, templateID, buildID, info); err != nil {
			buildFailed("%v", err)
			return
		}
	} else {
		buildParams := sandbox.StartTemplateBuildParams{}
		if info.FromImage != "" {
			buildParams.FromImage = &info.FromImage
		}
		if info.FromTemplate != "" {
			buildParams.FromTemplate = &info.FromTemplate
		}
		if info.StartCmd != "" {
			buildParams.StartCmd = &info.StartCmd
		}
		if info.ReadyCmd != "" {
			buildParams.ReadyCmd = &info.ReadyCmd
		}
		if info.NoCache {
			force := true
			buildParams.Force = &force
		}

		fmt.Printf("Starting build for template %s (build ID: %s)...\n", templateID, buildID)
		if err := client.StartTemplateBuild(ctx, templateID, buildID, buildParams); err != nil {
			buildFailed("start build failed: %v", err)
			return
		}
	}

	if !info.Wait {
		fmt.Printf("Build started. Use 'qshell sandbox template builds %s %s' to check status.\n", templateID, buildID)
		return
	}

	// 流式输出构建日志，支持 Ctrl+C 中断
	fmt.Println("Waiting for build to complete...")

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var cursor *int64
	for {
		logs, blErr := client.GetTemplateBuildLogs(ctx, templateID, buildID, &sandbox.GetBuildLogsParams{
			Cursor: cursor,
		})
		if blErr == nil && logs != nil {
			for _, entry := range logs.Logs {
				fmt.Printf("[%s] %s %s\n",
					sbClient.FormatTimestamp(entry.Timestamp),
					sbClient.LogLevelBadge(string(entry.Level)),
					entry.Message,
				)
				ts := entry.Timestamp.UnixMilli() + 1
				cursor = &ts
			}
		}

		// 检查构建状态
		buildInfo, bErr := client.GetTemplateBuildStatus(ctx, templateID, buildID, nil)
		if bErr != nil {
			buildFailed("get build status failed: %v", bErr)
			return
		}

		if buildInfo.Status == "ready" || buildInfo.Status == "error" {
			if buildInfo.Status == "error" {
				buildFailed("build failed")
				fmt.Print(formatBuildFailureLogs(buildInfo.Logs))
			} else {
				sbClient.PrintSuccess("Build completed!")
			}
			fmt.Printf("Template ID:  %s\n", buildInfo.TemplateID)
			fmt.Printf("Build ID:     %s\n", buildInfo.BuildID)
			fmt.Printf("Status:       %s\n", buildInfo.Status)

			if buildInfo.Status == "ready" {
				printSDKExamples(buildInfo.TemplateID)
			}
			return
		}

		select {
		case <-ctx.Done():
			buildFailed("build watch cancelled")
			return
		case <-time.After(3 * time.Second):
		}
	}
}

// cliBaseSelection 返回 CLI 显式指定的底模板字段与值。
func cliBaseSelection(info BuildInfo) (field, value string, err error) {
	var selected []string
	if info.TemplateID != "" {
		selected = append(selected, "--template-id")
	}
	if info.FromTemplate != "" {
		selected = append(selected, "--from-template")
	}
	if info.FromImage != "" {
		selected = append(selected, "--from-image")
	}
	if len(selected) > 1 {
		return "", "", fmt.Errorf("cannot use more than one of %s", strings.Join(selected, ", "))
	}
	switch {
	case info.TemplateID != "":
		return "template_id", info.TemplateID, nil
	case info.FromTemplate != "":
		return "from_template", info.FromTemplate, nil
	case info.FromImage != "":
		return "from_image", info.FromImage, nil
	}
	return "", "", nil
}

// applyCLIBasePrecedence 用 CLI 的底模板覆盖配置文件中的底模板字段，
// 使 --template-id / --from-template / --from-image 优先于 TOML。
// 返回是否覆盖了来自文件的值。
func applyCLIBasePrecedence(info *BuildInfo, field, value string) bool {
	overrode := (field != "template_id" && info.TemplateID != "") ||
		(field != "from_template" && info.FromTemplate != "") ||
		(field != "from_image" && info.FromImage != "")
	info.TemplateID, info.FromTemplate, info.FromImage = "", "", ""
	switch field {
	case "template_id":
		info.TemplateID = value
	case "from_template":
		info.FromTemplate = value
	case "from_image":
		info.FromImage = value
	}
	return overrode
}

// rejectSelfBase 拒绝把模板自己当作底模板。
//
// 旧版 qshell 会把新建模板的 ID 回写到 template_id；在新语义下该字段表示
// 底模板，于是构建会以模板自身为底，Dockerfile 的 FROM 被忽略。
func rejectSelfBase(info BuildInfo, existingID string, cfg *config.FileConfig, cliBase string) error {
	if info.FromTemplate == "" || info.FromTemplate != existingID {
		return nil
	}
	if cliBase != "" {
		flag := "--" + strings.ReplaceAll(cliBase, "_", "-")
		return fmt.Errorf("%s names the template being built (%s); choose a different base template", flag, existingID)
	}
	source := "qshell.sandbox.toml"
	if cfg != nil && cfg.SourcePath() != "" {
		source = cfg.SourcePath()
	}
	return fmt.Errorf("%s: template_id/from_template points at the template being built (%s); "+
		"remove that line so the build uses its declared base", source, existingID)
}

// baseConflictError 返回配置文件里同时声明多个底模板时的错误。
//
// 命令行显式指定底模板时以命令行为准，配置文件里的旧值由
// applyCLIBasePrecedence 清理；只有当冲突全部来自配置文件时才报错。
func baseConflictError(cfg *config.FileConfig, cliBase string) error {
	if cfg == nil || cliBase != "" {
		return nil
	}
	var keys []string
	if cfg.TemplateID != "" {
		keys = append(keys, "template_id")
	}
	if cfg.FromTemplate != "" {
		keys = append(keys, "from_template")
	}
	if cfg.FromImage != "" {
		keys = append(keys, "from_image")
	}
	if len(keys) < 2 {
		return nil
	}
	return fmt.Errorf("%s: %s all describe the base template; keep one and delete any template_id "+
		"line written by an older qshell", cfg.SourcePath(), strings.Join(keys, ", "))
}

// resolveBaseTemplate 合并 --template-id 与 --from-template：两者都表示
// 构建所用的底模板，因此不允许同时指定。
func resolveBaseTemplate(templateID, fromTemplate string) (string, error) {
	if templateID != "" && fromTemplate != "" {
		return "", fmt.Errorf("cannot use both --template-id and --from-template as the base template")
	}
	if templateID != "" {
		return templateID, nil
	}
	return fromTemplate, nil
}

func formatBuildFailureLogs(logs []string) string {
	if len(logs) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\nBuild Logs:\n")
	for _, log := range logs {
		fmt.Fprintf(&builder, "  %s\n", log)
	}
	return builder.String()
}

// ensureTemplateOwned 校验按 name 命中的模板是否属于当前团队。
//
// alias 解析可能命中其他团队的公共模板：此时 v3 创建会落在自己团队下，
// 但后续按 name 解析仍指向那个公共模板，构建结果不可达。提前失败比
// 「构建成功但沙箱用旧模板」更容易排查。
func ensureTemplateOwned(ctx context.Context, client *sandbox.Client, name, templateID string) error {
	tmpl, err := client.GetTemplate(ctx, templateID, nil)
	if err != nil {
		return fmt.Errorf("verify template %q ownership: %w", name, err)
	}
	if tmpl == nil {
		return nil
	}
	return ownershipError(name, templateID, tmpl.IsOwner)
}

// ownershipError 返回模板归属校验失败时的错误；属于当前团队时返回 nil。
func ownershipError(name, templateID string, isOwner bool) error {
	if isOwner {
		return nil
	}
	return fmt.Errorf("template %q (%s) belongs to another team and is only visible as public; "+
		"choose a unique template name owned by your team", name, templateID)
}

// legacyDiskSizeWarning 返回 qshell.sandbox.toml 中已废弃 disk_size_mb 的迁移提示。
//
// go-sdk v7.30.0 起，创建模板请求移除了 diskSizeMB，改为 minFreeDiskMb
// （构建后的最小剩余空间）。若继续接受旧字段却不提示，配置会静默失效，
// 用户会得到一个实际仍使用团队默认磁盘的模板。
func legacyDiskSizeWarning(cfg *config.FileConfig) string {
	if cfg == nil || cfg.DiskSizeMB == nil {
		return ""
	}
	return fmt.Sprintf("disk_size_mb = %d in %s is no longer supported: the template API removed the disk-size field, "+
		"so this setting has no effect; use min_free_disk_mb (minimum free space after build in MiB) instead",
		*cfg.DiskSizeMB, cfg.SourcePath())
}

func createTemplateParams(info BuildInfo) sandbox.CreateTemplateParams {
	params := sandbox.CreateTemplateParams{
		Name: &info.Name,
	}
	if info.CPUCount > 0 {
		params.CPUCount = &info.CPUCount
	}
	if info.MemoryMB > 0 {
		params.MemoryMB = &info.MemoryMB
	}
	if info.MinFreeDiskMB != nil {
		params.MinFreeDiskMB = info.MinFreeDiskMB
	}
	return params
}

// buildFromDockerfile 处理 v2 Dockerfile 构建流程：
// 解析 Dockerfile → 上传 COPY 文件 → 使用 steps 启动构建。
func buildFromDockerfile(ctx context.Context, client *sandbox.Client, templateID, buildID string, info BuildInfo) error {
	// 读取 Dockerfile
	content, err := os.ReadFile(info.Dockerfile)
	if err != nil {
		return fmt.Errorf("read Dockerfile: %w", err)
	}

	// 确定构建上下文目录
	contextPath := info.Path
	if contextPath == "" {
		contextPath = filepath.Dir(info.Dockerfile)
	}
	contextPath, err = filepath.Abs(contextPath)
	if err != nil {
		return fmt.Errorf("resolve context path: %w", err)
	}

	// 解析 Dockerfile
	result, err := dockerfile.Convert(string(content))
	if err != nil {
		return fmt.Errorf("parse Dockerfile: %w", err)
	}
	fmt.Printf("Parsed Dockerfile: base image=%s, %d steps\n", result.BaseImage, len(result.Steps))

	// 读取 .dockerignore
	ignorePatterns := dockerfile.ReadDockerignore(contextPath)

	// 处理 COPY 步骤：计算文件哈希并上传文件
	for i := range result.Steps {
		step := &result.Steps[i]
		if step.Type != "COPY" || step.Args == nil || len(*step.Args) < 2 {
			continue
		}
		args := *step.Args
		src, dest := args[0], args[1]

		// 计算文件哈希
		hash, err := dockerfile.ComputeFilesHash(src, dest, contextPath, ignorePatterns)
		if err != nil {
			return fmt.Errorf("compute file hash for COPY %s %s: %w", src, dest, err)
		}
		step.FilesHash = &hash

		// 检查文件是否需要上传
		fileInfo, err := client.GetTemplateFiles(ctx, templateID, hash)
		if err != nil {
			return fmt.Errorf("get template files for hash %s: %w", hash, err)
		}

		if !fileInfo.Present && fileInfo.URL != nil {
			fmt.Printf("Uploading files for COPY %s %s...\n", src, dest)
			if err := dockerfile.CollectAndUpload(ctx, *fileInfo.URL, src, contextPath, ignorePatterns); err != nil {
				return fmt.Errorf("upload files for COPY %s %s: %w", src, dest, err)
			}
		} else if fileInfo.Present {
			fmt.Printf("Files for COPY %s %s already uploaded (cached)\n", src, dest)
		}
	}

	// 构建参数
	buildParams := buildParamsFromDockerfileResult(result, info)

	// 应用来自 Dockerfile 或 CLI 覆盖的启动/就绪命令
	startCmd := result.StartCmd
	if info.StartCmd != "" {
		startCmd = info.StartCmd
	}
	if startCmd != "" {
		buildParams.StartCmd = &startCmd
	}

	readyCmd := result.ReadyCmd
	if info.ReadyCmd != "" {
		readyCmd = info.ReadyCmd
	}
	if readyCmd != "" {
		buildParams.ReadyCmd = &readyCmd
	}

	if info.NoCache {
		force := true
		buildParams.Force = &force
	}

	fmt.Printf("Starting build for template %s (build ID: %s)...\n", templateID, buildID)
	if err := client.StartTemplateBuild(ctx, templateID, buildID, buildParams); err != nil {
		return fmt.Errorf("start build: %w", err)
	}

	return nil
}

func validateBuildSourceSelection(info BuildInfo) error {
	if info.FromImage != "" && info.FromTemplate != "" {
		return fmt.Errorf("cannot use both --from-image and --from-template as the base template")
	}
	if info.Dockerfile == "" && info.FromImage == "" && info.FromTemplate == "" {
		return fmt.Errorf("--from-image, --from-template, or --dockerfile is required")
	}
	return nil
}

func buildParamsFromDockerfileResult(result *dockerfile.ConvertResult, info BuildInfo) sandbox.StartTemplateBuildParams {
	buildParams := sandbox.StartTemplateBuildParams{
		Steps: &result.Steps,
	}

	switch {
	case info.FromTemplate != "":
		buildParams.FromTemplate = &info.FromTemplate
	case info.FromImage != "":
		buildParams.FromImage = &info.FromImage
	default:
		buildParams.FromImage = &result.BaseImage
	}

	return buildParams
}

// printSDKExamples prints SDK usage examples for the given template ID.
func printSDKExamples(templateID string) {
	fmt.Println()
	sbClient.PrintSuccessBox("Template is ready! Use it with the SDK:")

	fmt.Printf("\n%s\n", sbClient.ColorInfo.Sprint("Go:"))
	fmt.Println(sbClient.FormatCodeBlock(fmt.Sprintf(`sb, _ := client.CreateAndWait(ctx, sandbox.CreateParams{
    TemplateID: "%s",
})`, templateID), "go"))

	fmt.Printf("\n%s\n", sbClient.ColorInfo.Sprint("Python:"))
	fmt.Println(sbClient.FormatCodeBlock(fmt.Sprintf(`sandbox = client.sandboxes.create("%s")`, templateID), "python"))

	fmt.Printf("\n%s\n", sbClient.ColorInfo.Sprint("TypeScript:"))
	fmt.Println(sbClient.FormatCodeBlock(fmt.Sprintf(`const sandbox = await client.sandboxes.create("%s")`, templateID), "typescript"))
	fmt.Println()
}
