package operations

import (
	"context"
	"fmt"
	"net/url"
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
	typ := strings.ToLower(strings.TrimSpace(fields["type"]))
	apiKey := optionalInlineString(fields["api-key"])
	baseURL := strings.TrimSpace(fields["base-url"])

	switch typ {
	case "openai":
		return sandbox.SandboxInjectionSpec{
			OpenAI: &sandbox.OpenAIInjection{
				APIKey:  apiKey,
				BaseURL: optionalInlineString(baseURL),
			},
		}, nil
	case "anthropic":
		return sandbox.SandboxInjectionSpec{
			Anthropic: &sandbox.AnthropicInjection{
				APIKey:  apiKey,
				BaseURL: optionalInlineString(baseURL),
			},
		}, nil
	case "gemini":
		return sandbox.SandboxInjectionSpec{
			Gemini: &sandbox.GeminiInjection{
				APIKey:  apiKey,
				BaseURL: optionalInlineString(baseURL),
			},
		}, nil
	case "http":
		if baseURL == "" {
			return sandbox.SandboxInjectionSpec{}, fmt.Errorf("inline injection type=http requires base-url")
		}
		parsedURL, err := url.Parse(baseURL)
		if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			return sandbox.SandboxInjectionSpec{}, fmt.Errorf("inline injection base-url must be a valid http/https URL")
		}
		headers := parseInlineHeaders(fields["headers"])
		httpInjection := &sandbox.HTTPInjection{
			BaseURL: baseURL,
		}
		if len(headers) > 0 {
			httpInjection.Headers = &headers
		}
		return sandbox.SandboxInjectionSpec{HTTP: httpInjection}, nil
	case "":
		return sandbox.SandboxInjectionSpec{}, fmt.Errorf("inline injection spec requires type")
	default:
		return sandbox.SandboxInjectionSpec{}, fmt.Errorf("unsupported inline injection type %q", typ)
	}
}

func parseInlineInjectionFields(spec string) map[string]string {
	fields := make(map[string]string)
	for _, part := range strings.Split(spec, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		fields[key] = value
	}
	return fields
}

func parseInlineHeaders(raw string) map[string]string {
	headers := make(map[string]string)
	for _, part := range strings.Split(raw, ";") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		headers[key] = value
	}
	return headers
}

func optionalInlineString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
