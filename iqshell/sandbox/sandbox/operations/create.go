package operations

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// CreateInfo holds parameters for creating a sandbox.
type CreateInfo struct {
	TemplateID      string
	Timeout         int32
	Metadata        string
	Detach          bool
	EnvVars         []string // KEY=VALUE pairs
	AutoPause       bool
	InjectionRuleID []string
	InlineInjection []string
	// Resources 沙箱启动前挂载的资源规约（如 GitHub 仓库、Kodo bucket），格式参见 parseSandboxResource
	Resources []string
}

// Create creates a new sandbox and connects to its terminal.
// When the terminal session ends, the sandbox is killed.
// The sandbox stays alive via keep-alive in the terminal session (matches e2b CLI behavior).
func Create(info CreateInfo) {
	if info.TemplateID == "" {
		sbClient.PrintError("template ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	params := sandbox.CreateParams{
		TemplateID: info.TemplateID,
	}
	if info.Timeout > 0 {
		params.Timeout = &info.Timeout
	}
	if info.Metadata != "" {
		meta := sandbox.Metadata(sbClient.ParseMetadataMap(info.Metadata))
		params.Metadata = &meta
	}
	if len(info.EnvVars) > 0 {
		envMap := parseEnvPairs(info.EnvVars)
		if len(envMap) > 0 {
			params.EnvVars = &envMap
		}
	}
	if info.AutoPause {
		params.AutoPause = &info.AutoPause
	}
	injections, err := buildSandboxInjections(info.InjectionRuleID, info.InlineInjection)
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}
	if len(injections) > 0 {
		params.Injections = &injections
	}
	resources, err := buildSandboxResources(info.Resources)
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}
	if len(resources) > 0 {
		params.Resources = &resources
	}

	fmt.Printf("Creating sandbox from template %s...\n", info.TemplateID)
	sb, _, err := client.CreateAndWait(ctx, params)
	if err != nil {
		sbClient.PrintError("create sandbox failed: %v", err)
		return
	}
	if info.Detach {
		sbClient.PrintSuccess("Sandbox %s created", sb.ID())
		fmt.Printf("Sandbox ID:   %s\n", sb.ID())
		fmt.Printf("Template ID:  %s\n", sb.TemplateID())
		fmt.Println()
		fmt.Printf("Connect:  qshell sandbox connect %s\n", sb.ID())
		fmt.Printf("Exec:     qshell sandbox exec %s -- <command>\n", sb.ID())
		fmt.Printf("Kill:     qshell sandbox kill %s\n", sb.ID())
		return
	}

	sbClient.PrintSuccess("Sandbox %s created, connecting...", sb.ID())

	// When create session ends, kill the sandbox
	defer func() {
		fmt.Printf("\nKilling sandbox %s...\n", sb.ID())
		if kErr := sb.Kill(context.Background()); kErr != nil {
			// Ignore 404 errors: sandbox may have already been terminated by timeout
			if !strings.Contains(kErr.Error(), "404") {
				sbClient.PrintWarn("kill sandbox failed: %v", kErr)
			}
		}
	}()

	runTerminalSession(ctx, sb)
}

func buildSandboxInjections(ruleIDs, inlineSpecs []string) ([]sandbox.SandboxInjectionSpec, error) {
	if len(ruleIDs) == 0 && len(inlineSpecs) == 0 {
		return nil, nil
	}

	injections := make([]sandbox.SandboxInjectionSpec, 0, len(ruleIDs)+len(inlineSpecs))
	for _, ruleID := range ruleIDs {
		trimmed := strings.TrimSpace(ruleID)
		if trimmed == "" {
			return nil, fmt.Errorf("injection rule ID cannot be empty")
		}
		injections = append(injections, sandbox.SandboxInjectionSpec{
			ByID: &trimmed,
		})
	}
	for _, spec := range inlineSpecs {
		injection, err := parseInlineSandboxInjection(spec)
		if err != nil {
			return nil, err
		}
		injections = append(injections, injection)
	}
	return injections, nil
}

func parseInlineSandboxInjection(spec string) (sandbox.SandboxInjectionSpec, error) {
	fields := parseInlineInjectionFields(spec)
	parts, err := sbClient.BuildInjectionParts(
		fields["type"],
		fields["api-key"],
		fields["base-url"],
		parseInlineHeaders(fields["headers"]),
		parseInlineHeaders(fields["if-headers"]),
		parseInlineHeaders(fields["if-queries"]),
	)
	if err != nil {
		return sandbox.SandboxInjectionSpec{}, fmt.Errorf("invalid inline injection spec: %w", err)
	}
	return sandbox.SandboxInjectionSpec{
		OpenAI:    parts.OpenAI,
		Anthropic: parts.Anthropic,
		Gemini:    parts.Gemini,
		Qiniu:     parts.Qiniu,
		Github:    parts.Github,
		HTTP:      parts.HTTP,
	}, nil
}

func parseInlineInjectionFields(spec string) map[string]string {
	fields := make(map[string]string)
	currentKey := ""
	for _, part := range strings.Split(spec, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" && isInlineInjectionMapField(currentKey) && fields[currentKey] != "" {
				fields[currentKey] += "," + trimmed
			}
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		if !isInlineInjectionFieldKey(key) {
			if isInlineInjectionMapField(currentKey) {
				fields[currentKey] += "," + strings.TrimSpace(part)
			}
			continue
		}
		fields[key] = value
		currentKey = key
	}
	return fields
}

func isInlineInjectionMapField(key string) bool {
	switch key {
	case "headers", "if-headers", "if-queries":
		return true
	default:
		return false
	}
}

func isInlineInjectionFieldKey(key string) bool {
	switch key {
	case "type", "api-key", "base-url", "headers", "if-headers", "if-queries":
		return true
	default:
		return false
	}
}

// buildSandboxResources 把命令行传入的 --resource 规约转换为 SDK 的资源列表。
func buildSandboxResources(resourceSpecs []string) ([]sandbox.SandboxResourceSpec, error) {
	if len(resourceSpecs) == 0 {
		return nil, nil
	}
	resources := make([]sandbox.SandboxResourceSpec, 0, len(resourceSpecs))
	// 同一沙箱内多个 GitHub 仓库资源当前必须共用同一 token（go-sdk 注释明示约束）；
	// 提前在 CLI 层校验，避免等到平台克隆阶段才返回不易理解的错误。
	var seenToken string
	for _, spec := range resourceSpecs {
		resource, err := parseSandboxResource(spec)
		if err != nil {
			return nil, err
		}
		if gr := resource.GitRepository; gr != nil {
			if gr.AuthorizationToken == nil {
				return nil, fmt.Errorf("invalid resource spec %q: token is required for github_repository", spec)
			}
			switch token := *gr.AuthorizationToken; {
			case seenToken == "":
				seenToken = token
			case token != seenToken:
				return nil, fmt.Errorf("inconsistent --resource tokens: a sandbox can carry only one GitHub token across all repository resources")
			}
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

// parseSandboxResource 解析单条 --resource 规约。
// 支持格式：
//   - type=github_repository,url=<url>,mount-path=<absPath>,token=<token>
//   - type=kodo,bucket=<bucket>,mount-path=<absPath>[,prefix=<prefix>][,read-only=<bool>]
func parseSandboxResource(spec string) (sandbox.SandboxResourceSpec, error) {
	fields := sbClient.ParseMetadataMap(spec)

	typ := strings.ToLower(fields["type"])
	if typ == "" {
		typ = string(sandbox.GitRepositoryTypeGithub)
	}

	switch typ {
	case string(sandbox.GitRepositoryTypeGithub):
		url := fields["url"]
		if url == "" {
			return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: url is required for github_repository", spec)
		}
		mountPath, err := parseSandboxResourceMountPath(spec, fields, "github_repository")
		if err != nil {
			return sandbox.SandboxResourceSpec{}, err
		}
		token := fields["token"]
		if token == "" {
			return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: token is required for github_repository", spec)
		}
		res := &sandbox.GitRepositoryResource{
			Type:               sandbox.GitRepositoryTypeGithub,
			URL:                url,
			MountPath:          mountPath,
			AuthorizationToken: &token,
		}
		return sandbox.SandboxResourceSpec{GitRepository: res}, nil
	case "kodo":
		bucket := fields["bucket"]
		if bucket == "" {
			return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: bucket is required for kodo", spec)
		}
		mountPath, err := parseSandboxResourceMountPath(spec, fields, "kodo")
		if err != nil {
			return sandbox.SandboxResourceSpec{}, err
		}
		res := &sandbox.KodoResource{
			Bucket:    bucket,
			MountPath: mountPath,
		}
		if prefix := fields["prefix"]; prefix != "" {
			res.Prefix = &prefix
		}
		readOnly := fields["read-only"]
		readOnlyAlias := fields["readonly"]
		if readOnly != "" && readOnlyAlias != "" && readOnly != readOnlyAlias {
			return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: read-only %q and readonly %q conflict, specify only one", spec, readOnly, readOnlyAlias)
		}
		if readOnly == "" {
			readOnly = readOnlyAlias
		}
		if readOnly != "" {
			value, err := strconv.ParseBool(readOnly)
			if err != nil {
				return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: read-only %q must be a boolean", spec, readOnly)
			}
			res.ReadOnly = &value
		}
		return sandbox.SandboxResourceSpec{Kodo: res}, nil
	default:
		return sandbox.SandboxResourceSpec{}, fmt.Errorf("invalid resource spec %q: unsupported type %q (supported: github_repository, kodo)", spec, typ)
	}
}

func parseSandboxResourceMountPath(spec string, fields map[string]string, resourceType string) (string, error) {
	mountPath := fields["mount-path"]
	mountAlias := fields["mount"]
	// 同时给出 mount-path= 与 mount= 且取值不一致时直接报错，避免静默忽略其中一项造成误解
	if mountPath != "" && mountAlias != "" && mountPath != mountAlias {
		return "", fmt.Errorf("invalid resource spec %q: mount-path %q and mount %q conflict, specify only one", spec, mountPath, mountAlias)
	}
	if mountPath == "" {
		// 兼容 mount= 简写
		mountPath = mountAlias
	}
	if mountPath == "" {
		return "", fmt.Errorf("invalid resource spec %q: mount-path is required for %s", spec, resourceType)
	}
	// 沙箱内部使用 POSIX 路径；用 path.IsAbs 而非 filepath.IsAbs，避免 Windows 主机上把 /workspace 误判为相对
	if !path.IsAbs(mountPath) {
		return "", fmt.Errorf("invalid resource spec %q: mount-path %q must be an absolute path", spec, mountPath)
	}
	return mountPath, nil
}

func parseInlineHeaders(raw string) map[string]string {
	m := make(map[string]string)
	if raw == "" {
		return m
	}
	separator := ";"
	if !strings.Contains(raw, separator) {
		separator = ","
	}
	for _, pair := range strings.Split(raw, separator) {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		key, value, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		m[key] = strings.TrimSpace(value)
	}
	return m
}
