package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/config"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/dockerfile"
)

// BuildInfo 保存构建模板的参数。
type BuildInfo struct {
	// Name 是模板名称（用于创建 + 构建）。
	Name string

	// TemplateID 是已有模板 ID（用于重新构建）。
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

// Build 创建或重新构建模板。
// 如果提供了 TemplateID，则对已有模板启动新构建。
// 否则，使用给定的 Name 创建新模板并启动构建。
func Build(info BuildInfo) {
	// 在合并配置前捕获"CLI 未提供 TemplateID"，以便后续判断是否需要回写。
	noIDBeforeMerge := info.TemplateID == ""
	cliFromImage := info.FromImage != ""
	cliFromTemplate := info.FromTemplate != ""

	// 加载配置文件并合并（CLI > file > default）
	var cfg *config.FileConfig
	if info.ConfigPath != "" {
		loaded, cErr := config.Load(info.ConfigPath)
		if cErr != nil {
			sbClient.PrintError("load config: %v", cErr)
			return
		}
		if loaded == nil {
			sbClient.PrintError("config file not found: %s", info.ConfigPath)
			return
		}
		cfg = loaded
	} else {
		loaded, cErr := config.LoadFromCwd()
		if cErr != nil {
			sbClient.PrintError("load config: %v", cErr)
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
		info.NoCache = fields.NoCache

		for _, key := range overrides {
			fmt.Fprintf(os.Stderr, "[config] CLI overrides %s from %s\n", key, cfg.SourcePath())
		}
	}

	if err := normalizeRebuildSourceSelection(&info, cliFromImage, cliFromTemplate); err != nil {
		sbClient.PrintError("%v", err)
		return
	}
	if err := validateBuildSourceSelection(info); err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	templateID := info.TemplateID
	buildID := ""

	if templateID == "" {
		// 创建新模板
		if info.Name == "" {
			sbClient.PrintError("template name (--name) or template ID (--template-id) is required")
			return
		}

		createParams := sandbox.CreateTemplateParams{
			Name: &info.Name,
		}
		if info.CPUCount > 0 {
			createParams.CPUCount = &info.CPUCount
		}
		if info.MemoryMB > 0 {
			createParams.MemoryMB = &info.MemoryMB
		}

		fmt.Printf("Creating template %s...\n", info.Name)
		resp, cErr := client.CreateTemplate(ctx, createParams)
		if cErr != nil {
			sbClient.PrintError("create template failed: %v", cErr)
			return
		}
		templateID = resp.TemplateID
		buildID = resp.BuildID
		sbClient.PrintSuccess("Template %s created (build ID: %s)", templateID, buildID)
	} else {
		// 对已有模板申请新的 waiting build（对齐 E2B CLI 的 rebuild 流程）：
		// RebuildTemplate → StartTemplateBuild → WaitForBuild。
		// 直接复用历史 build ID 会导致 400 "build is not in waiting state"。
		dockerfileContent, dErr := dockerfileForRebuild(info)
		if dErr != nil {
			sbClient.PrintError("%v", dErr)
			return
		}
		rebuildParams := sandbox.RebuildTemplateParams{
			Dockerfile: dockerfileContent,
		}
		if info.CPUCount > 0 {
			rebuildParams.CPUCount = &info.CPUCount
		}
		if info.MemoryMB > 0 {
			rebuildParams.MemoryMB = &info.MemoryMB
		}
		if info.StartCmd != "" {
			rebuildParams.StartCmd = &info.StartCmd
		}
		if info.ReadyCmd != "" {
			rebuildParams.ReadyCmd = &info.ReadyCmd
		}

		fmt.Printf("Requesting new build for template %s...\n", templateID)
		resp, rErr := client.RebuildTemplate(ctx, templateID, rebuildParams)
		if rErr != nil {
			sbClient.PrintError("rebuild template failed: %v", rErr)
			return
		}
		buildID = resp.BuildID
		sbClient.PrintSuccess("New build requested (build ID: %s)", buildID)
	}

	if info.Dockerfile != "" {
		if err := buildFromDockerfile(ctx, client, templateID, buildID, info); err != nil {
			sbClient.PrintError("%v", err)
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
			sbClient.PrintError("start build failed: %v", err)
			return
		}
	}

	if !info.Wait {
		writeTemplateIDToConfigIfNeeded(cfg, noIDBeforeMerge, templateID)
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
			sbClient.PrintError("get build status failed: %v", bErr)
			return
		}

		if buildInfo.Status == "ready" || buildInfo.Status == "error" {
			if buildInfo.Status == "error" {
				sbClient.PrintError("build failed")
			} else {
				sbClient.PrintSuccess("Build completed!")
				writeTemplateIDToConfigIfNeeded(cfg, noIDBeforeMerge, templateID)
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
			sbClient.PrintError("build watch cancelled")
			return
		case <-time.After(3 * time.Second):
		}
	}
}

func writeTemplateIDToConfigIfNeeded(cfg *config.FileConfig, noIDBeforeMerge bool, templateID string) {
	if cfg == nil || !noIDBeforeMerge || cfg.SourcePath() == "" {
		return
	}
	if wErr := config.WriteTemplateID(cfg.SourcePath(), templateID); wErr != nil {
		fmt.Fprintf(os.Stderr, "[config] warning: failed to write template_id to %s: %v\n",
			cfg.SourcePath(), wErr)
		return
	}
	sbClient.PrintSuccess("Written template_id to %s (please commit this file)", cfg.SourcePath())
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
		return fmt.Errorf("cannot specify both --from-image and --from-template")
	}
	if info.Dockerfile == "" && info.FromImage == "" && info.FromTemplate == "" {
		return fmt.Errorf("--from-image, --from-template, or --dockerfile is required")
	}
	return nil
}

func normalizeRebuildSourceSelection(info *BuildInfo, cliFromImage, cliFromTemplate bool) error {
	if info.TemplateID == "" || (info.FromImage == "" && info.FromTemplate == "") {
		return nil
	}
	if cliFromImage || cliFromTemplate {
		return fmt.Errorf("cannot specify --from-image or --from-template when rebuilding an existing template (--template-id)")
	}
	info.FromImage = ""
	info.FromTemplate = ""
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

// dockerfileForRebuild 返回 rebuild 请求所需的 Dockerfile 文本。
// E2B v1 rebuild API（POST /templates/{id}）强制要求在请求体中携带
// Dockerfile 内容，因此 --template-id 场景必须同时提供 --dockerfile；
// --from-image / --from-template 仅适用于新建模板。
func dockerfileForRebuild(info BuildInfo) (string, error) {
	if info.Dockerfile == "" {
		return "", fmt.Errorf("--dockerfile is required when rebuilding an existing template (--template-id); --from-image and --from-template only apply to new templates")
	}
	content, err := os.ReadFile(info.Dockerfile)
	if err != nil {
		return "", fmt.Errorf("read Dockerfile: %w", err)
	}
	return string(content), nil
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
