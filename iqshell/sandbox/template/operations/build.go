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

	// Dockerfile 是 Dockerfile 的路径（启用 v2 Dockerfile 构建）。
	Dockerfile string

	// Path 是构建上下文目录（默认为 Dockerfile 所在目录）。
	Path string
}

// Build 创建或重新构建模板。
// 如果提供了 TemplateID，则对已有模板启动新构建。
// 否则，使用给定的 Name 创建新模板并启动构建。
func Build(info BuildInfo) {
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
		// 获取已有模板以查找最新的 build ID
		tmpl, gErr := client.GetTemplate(ctx, templateID, nil)
		if gErr != nil {
			sbClient.PrintError("get template failed: %v", gErr)
			return
		}
		if len(tmpl.Builds) > 0 {
			// 使用最后一个 build（API 按时间升序返回，最新的在末尾）
			buildID = tmpl.Builds[len(tmpl.Builds)-1].BuildID
		} else {
			sbClient.PrintError("no builds found for template, cannot rebuild")
			return
		}
	}

	if info.Dockerfile != "" {
		if err := buildFromDockerfile(ctx, client, templateID, buildID, info); err != nil {
			sbClient.PrintError("%v", err)
			return
		}
	} else {
		// 验证构建来源
		if info.FromImage == "" && info.FromTemplate == "" {
			sbClient.PrintError("--from-image, --from-template, or --dockerfile is required")
			return
		}

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
	buildParams := sandbox.StartTemplateBuildParams{
		FromImage: &result.BaseImage,
		Steps:     &result.Steps,
	}

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
