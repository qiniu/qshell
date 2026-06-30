package operations

import (
	"slices"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

const (
	injectionTypeOpenAI    = "openai"
	injectionTypeAnthropic = "anthropic"
	injectionTypeGemini    = "gemini"
	injectionTypeQiniu     = "qiniu"
	injectionTypeGithub    = "github"
	injectionTypeHTTP      = "http"
)

// githubInjectionTarget 描述 GitHub 注入由平台侧匹配的目标域名集合。
// 当前平台行为固定为 github.com / api.github.com；若后续 SDK 支持自托管 GitHub Enterprise，
// 需要随 SDK 暴露的常量一起同步更新此处。
const githubInjectionTarget = "github.com, api.github.com"

type injectionInput struct {
	Type      string
	APIKey    string
	BaseURL   string
	Headers   string
	IfHeaders string
	IfQueries string
}

func buildInjectionSpec(input injectionInput) (sandbox.InjectionSpec, error) {
	parts, err := sbClient.BuildInjectionParts(
		input.Type,
		input.APIKey,
		input.BaseURL,
		sbClient.ParseMetadataMap(input.Headers),
		sbClient.ParseMetadataMap(input.IfHeaders),
		sbClient.ParseMetadataMap(input.IfQueries),
	)
	if err != nil {
		return sandbox.InjectionSpec{}, err
	}
	return sandbox.InjectionSpec{
		OpenAI:    parts.OpenAI,
		Anthropic: parts.Anthropic,
		Gemini:    parts.Gemini,
		Qiniu:     parts.Qiniu,
		Github:    parts.Github,
		HTTP:      parts.HTTP,
	}, nil
}

func shouldUpdateInjection(input injectionInput) bool {
	return strings.TrimSpace(input.Type) != "" ||
		strings.TrimSpace(input.APIKey) != "" ||
		strings.TrimSpace(input.BaseURL) != "" ||
		strings.TrimSpace(input.Headers) != "" ||
		strings.TrimSpace(input.IfHeaders) != "" ||
		strings.TrimSpace(input.IfQueries) != ""
}

func formatInjectionType(spec sandbox.InjectionSpec) string {
	switch {
	case spec.OpenAI != nil:
		return injectionTypeOpenAI
	case spec.Anthropic != nil:
		return injectionTypeAnthropic
	case spec.Gemini != nil:
		return injectionTypeGemini
	case spec.Qiniu != nil:
		return injectionTypeQiniu
	case spec.Github != nil:
		return injectionTypeGithub
	case spec.HTTP != nil:
		return injectionTypeHTTP
	default:
		return "-"
	}
}

func formatInjectionTarget(spec sandbox.InjectionSpec) string {
	switch {
	case spec.OpenAI != nil:
		return optionalValue(spec.OpenAI.BaseURL, "api.openai.com")
	case spec.Anthropic != nil:
		return optionalValue(spec.Anthropic.BaseURL, "api.anthropic.com")
	case spec.Gemini != nil:
		return optionalValue(spec.Gemini.BaseURL, "generativelanguage.googleapis.com")
	case spec.Qiniu != nil:
		return optionalValue(spec.Qiniu.BaseURL, "api.qnaigc.com")
	case spec.Github != nil:
		return optionalValue(spec.Github.BaseURL, githubInjectionTarget)
	case spec.HTTP != nil:
		return spec.HTTP.BaseURL
	default:
		return "-"
	}
}

func formatInjectionHeaders(spec sandbox.InjectionSpec) string {
	if spec.HTTP == nil || spec.HTTP.Headers == nil || len(*spec.HTTP.Headers) == 0 {
		return "-"
	}
	return formatMapKeys(*spec.HTTP.Headers)
}

func formatInjectionConditions(spec sandbox.InjectionSpec) string {
	headers, queries := injectionConditions(spec)
	parts := make([]string, 0, 2)
	if headers != nil && len(*headers) > 0 {
		parts = append(parts, "headers: "+formatMapKeys(*headers))
	}
	if queries != nil && len(*queries) > 0 {
		parts = append(parts, "queries: "+formatMapKeys(*queries))
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, "; ")
}

func injectionConditions(spec sandbox.InjectionSpec) (*map[string]string, *map[string]string) {
	switch {
	case spec.OpenAI != nil:
		return spec.OpenAI.IfHeaders, spec.OpenAI.IfQueries
	case spec.Anthropic != nil:
		return spec.Anthropic.IfHeaders, spec.Anthropic.IfQueries
	case spec.Gemini != nil:
		return spec.Gemini.IfHeaders, spec.Gemini.IfQueries
	case spec.Qiniu != nil:
		return spec.Qiniu.IfHeaders, spec.Qiniu.IfQueries
	case spec.Github != nil:
		return spec.Github.IfHeaders, spec.Github.IfQueries
	case spec.HTTP != nil:
		return spec.HTTP.IfHeaders, spec.HTTP.IfQueries
	default:
		return nil, nil
	}
}

func formatMapKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

func hasAPIKey(spec sandbox.InjectionSpec) bool {
	switch {
	case spec.OpenAI != nil:
		return spec.OpenAI.APIKey != nil && strings.TrimSpace(*spec.OpenAI.APIKey) != ""
	case spec.Anthropic != nil:
		return spec.Anthropic.APIKey != nil && strings.TrimSpace(*spec.Anthropic.APIKey) != ""
	case spec.Gemini != nil:
		return spec.Gemini.APIKey != nil && strings.TrimSpace(*spec.Gemini.APIKey) != ""
	case spec.Qiniu != nil:
		return spec.Qiniu.APIKey != nil && strings.TrimSpace(*spec.Qiniu.APIKey) != ""
	case spec.Github != nil:
		return spec.Github.Token != nil && strings.TrimSpace(*spec.Github.Token) != ""
	default:
		return false
	}
}

func optionalValue(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}
