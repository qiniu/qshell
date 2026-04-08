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
	injectionTypeHTTP      = "http"
)

type injectionInput struct {
	Type    string
	APIKey  string
	BaseURL string
	Headers string
}

func buildInjectionSpec(input injectionInput) (sandbox.InjectionSpec, error) {
	parts, err := sbClient.BuildInjectionParts(input.Type, input.APIKey, input.BaseURL, sbClient.ParseMetadataMap(input.Headers))
	if err != nil {
		return sandbox.InjectionSpec{}, err
	}
	return sandbox.InjectionSpec{
		OpenAI:    parts.OpenAI,
		Anthropic: parts.Anthropic,
		Gemini:    parts.Gemini,
		HTTP:      parts.HTTP,
	}, nil
}

func shouldUpdateInjection(input injectionInput) bool {
	return strings.TrimSpace(input.Type) != "" || strings.TrimSpace(input.APIKey) != "" || strings.TrimSpace(input.BaseURL) != "" || strings.TrimSpace(input.Headers) != ""
}

func formatInjectionType(spec sandbox.InjectionSpec) string {
	switch {
	case spec.OpenAI != nil:
		return injectionTypeOpenAI
	case spec.Anthropic != nil:
		return injectionTypeAnthropic
	case spec.Gemini != nil:
		return injectionTypeGemini
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
	keys := make([]string, 0, len(*spec.HTTP.Headers))
	for k := range *spec.HTTP.Headers {
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
