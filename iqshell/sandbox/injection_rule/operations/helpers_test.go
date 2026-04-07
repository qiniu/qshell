package operations

import (
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

func TestBuildInjectionSpecOpenAI(t *testing.T) {
	apiKey := "sk-test"
	baseURL := "https://api.openai-proxy.example.com/v1"

	spec, err := buildInjectionSpec(injectionInput{
		Type:    injectionTypeOpenAI,
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("buildInjectionSpec() error = %v", err)
	}
	if spec.OpenAI == nil {
		t.Fatal("expected OpenAI injection to be set")
	}
	if spec.OpenAI.APIKey == nil || *spec.OpenAI.APIKey != apiKey {
		t.Fatalf("OpenAI API key = %v, want %q", spec.OpenAI.APIKey, apiKey)
	}
	if spec.OpenAI.BaseURL == nil || *spec.OpenAI.BaseURL != baseURL {
		t.Fatalf("OpenAI base URL = %v, want %q", spec.OpenAI.BaseURL, baseURL)
	}
}

func TestBuildInjectionSpecHTTP(t *testing.T) {
	spec, err := buildInjectionSpec(injectionInput{
		Type:    injectionTypeHTTP,
		BaseURL: "https://api.example.com",
		Headers: "Authorization=Bearer token,X-Env=prod",
	})
	if err != nil {
		t.Fatalf("buildInjectionSpec() error = %v", err)
	}
	if spec.HTTP == nil {
		t.Fatal("expected HTTP injection to be set")
	}
	if spec.HTTP.BaseURL != "https://api.example.com" {
		t.Fatalf("HTTP base URL = %q, want %q", spec.HTTP.BaseURL, "https://api.example.com")
	}
	if spec.HTTP.Headers == nil || len(*spec.HTTP.Headers) != 2 {
		t.Fatalf("HTTP headers = %v, want 2 headers", spec.HTTP.Headers)
	}
}

func TestBuildInjectionSpecRejectsMissingType(t *testing.T) {
	if _, err := buildInjectionSpec(injectionInput{}); err == nil {
		t.Fatal("expected missing type to fail")
	}
}

func TestBuildInjectionSpecRejectsHTTPWithoutBaseURL(t *testing.T) {
	if _, err := buildInjectionSpec(injectionInput{Type: injectionTypeHTTP}); err == nil {
		t.Fatal("expected HTTP injection without base URL to fail")
	}
}

func TestBuildInjectionSpecRejectsHTTPWithInvalidScheme(t *testing.T) {
	if _, err := buildInjectionSpec(injectionInput{
		Type:    injectionTypeHTTP,
		BaseURL: "file:///tmp/secret",
	}); err == nil {
		t.Fatal("expected HTTP injection with invalid scheme to fail")
	}
}

func TestFormatInjectionSummaryOpenAI(t *testing.T) {
	spec := sandbox.InjectionSpec{
		OpenAI: &sandbox.OpenAIInjection{},
	}

	if got := formatInjectionType(spec); got != "openai" {
		t.Fatalf("formatInjectionType() = %q, want %q", got, "openai")
	}
	if got := formatInjectionTarget(spec); got != "api.openai.com" {
		t.Fatalf("formatInjectionTarget() = %q, want %q", got, "api.openai.com")
	}
	if got := formatInjectionHeaders(spec); got != "-" {
		t.Fatalf("formatInjectionHeaders() = %q, want %q", got, "-")
	}
}

func TestFormatInjectionSummaryHTTP(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token",
		"X-Trace":       "true",
	}
	spec := sandbox.InjectionSpec{
		HTTP: &sandbox.HTTPInjection{
			BaseURL: "https://api.example.com",
			Headers: &headers,
		},
	}

	if got := formatInjectionType(spec); got != "http" {
		t.Fatalf("formatInjectionType() = %q, want %q", got, "http")
	}
	if got := formatInjectionTarget(spec); got != "https://api.example.com" {
		t.Fatalf("formatInjectionTarget() = %q, want %q", got, "https://api.example.com")
	}
	got := formatInjectionHeaders(spec)
	if got != "Authorization, X-Trace" && got != "X-Trace, Authorization" {
		t.Fatalf("formatInjectionHeaders() = %q, want header key list", got)
	}
}
