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
		Type:      injectionTypeHTTP,
		BaseURL:   "https://api.example.com",
		Headers:   "Authorization=Bearer token,X-Env=prod",
		IfHeaders: "X-Scope=demo",
		IfQueries: "inject=true",
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
	if spec.HTTP.IfHeaders == nil || (*spec.HTTP.IfHeaders)["X-Scope"] != "demo" {
		t.Fatalf("HTTP if headers = %v, want X-Scope=demo", spec.HTTP.IfHeaders)
	}
	if spec.HTTP.IfQueries == nil || (*spec.HTTP.IfQueries)["inject"] != "true" {
		t.Fatalf("HTTP if queries = %v, want inject=true", spec.HTTP.IfQueries)
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
	ifHeaders := map[string]string{"X-Scope": "demo"}
	ifQueries := map[string]string{"inject": "true"}
	spec := sandbox.InjectionSpec{
		HTTP: &sandbox.HTTPInjection{
			BaseURL:   "https://api.example.com",
			Headers:   &headers,
			IfHeaders: &ifHeaders,
			IfQueries: &ifQueries,
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
	if got := formatInjectionConditions(spec); got != "headers: X-Scope; queries: inject" {
		t.Fatalf("formatInjectionConditions() = %q, want condition key list", got)
	}
}

func TestFormatInjectionSummaryQiniu(t *testing.T) {
	baseURL := " https://qiniu-proxy.example.com "
	spec := sandbox.InjectionSpec{
		Qiniu: &sandbox.QiniuInjection{BaseURL: &baseURL},
	}

	if got := formatInjectionType(spec); got != "qiniu" {
		t.Fatalf("formatInjectionType() = %q, want %q", got, "qiniu")
	}
	if got := formatInjectionTarget(spec); got != "https://qiniu-proxy.example.com" {
		t.Fatalf("formatInjectionTarget() = %q, want %q", got, "https://qiniu-proxy.example.com")
	}
}

func TestFormatInjectionTargetQiniuDefault(t *testing.T) {
	spec := sandbox.InjectionSpec{Qiniu: &sandbox.QiniuInjection{}}
	if got := formatInjectionTarget(spec); got != "api.qnaigc.com" {
		t.Fatalf("formatInjectionTarget() = %q, want %q", got, "api.qnaigc.com")
	}
}

func TestFormatInjectionSummaryGithub(t *testing.T) {
	token := "ghp-token"
	baseURL := "https://api.github.com/repos/qiniu/*"
	spec := sandbox.InjectionSpec{
		Github: &sandbox.GithubInjection{BaseURL: &baseURL, Token: &token},
	}

	if got := formatInjectionType(spec); got != "github" {
		t.Fatalf("formatInjectionType() = %q, want %q", got, "github")
	}
	if got := formatInjectionTarget(spec); got != "https://api.github.com/repos/qiniu/*" {
		t.Fatalf("formatInjectionTarget() = %q, want %q", got, "https://api.github.com/repos/qiniu/*")
	}
	if got := formatInjectionHeaders(spec); got != "-" {
		t.Fatalf("formatInjectionHeaders() = %q, want %q", got, "-")
	}
	if !hasAPIKey(spec) {
		t.Fatal("hasAPIKey() = false, want true for github with token")
	}
}

func TestBuildInjectionSpecGithub(t *testing.T) {
	spec, err := buildInjectionSpec(injectionInput{
		Type:      injectionTypeGithub,
		APIKey:    "ghp-token",
		BaseURL:   "https://api.github.com/repos/qiniu/*",
		IfHeaders: "X-GitHub-Api-Version=2022-11-28",
	})
	if err != nil {
		t.Fatalf("buildInjectionSpec(github) error = %v", err)
	}
	if spec.Github == nil {
		t.Fatal("expected Github injection to be set")
	}
	if spec.Github.Token == nil || *spec.Github.Token != "ghp-token" {
		t.Fatalf("Github token = %v, want %q", spec.Github.Token, "ghp-token")
	}
	if spec.Github.BaseURL == nil || *spec.Github.BaseURL != "https://api.github.com/repos/qiniu/*" {
		t.Fatalf("Github base URL = %v, want https://api.github.com/repos/qiniu/*", spec.Github.BaseURL)
	}
	if spec.Github.IfHeaders == nil || (*spec.Github.IfHeaders)["X-GitHub-Api-Version"] != "2022-11-28" {
		t.Fatalf("Github if headers = %v, want version header", spec.Github.IfHeaders)
	}
}

// 空 token 的 github 规则无意义，需要在 CLI 层提前拒绝
func TestBuildInjectionSpecGithubRequiresToken(t *testing.T) {
	if _, err := buildInjectionSpec(injectionInput{Type: injectionTypeGithub}); err == nil {
		t.Fatal("expected buildInjectionSpec(github) without api-key to fail")
	}
}

func TestShouldUpdateInjection(t *testing.T) {
	if shouldUpdateInjection(injectionInput{}) {
		t.Fatal("shouldUpdateInjection() = true, want false")
	}
	if !shouldUpdateInjection(injectionInput{Type: " qiniu "}) {
		t.Fatal("shouldUpdateInjection() = false, want true when type is set")
	}
	if !shouldUpdateInjection(injectionInput{APIKey: "sk"}) {
		t.Fatal("shouldUpdateInjection() = false, want true when api key is set")
	}
	if !shouldUpdateInjection(injectionInput{IfHeaders: "X-Scope=demo"}) {
		t.Fatal("shouldUpdateInjection() = false, want true when if headers are set")
	}
}

func TestHasAPIKey(t *testing.T) {
	key := "sk-qiniu"
	spec := sandbox.InjectionSpec{Qiniu: &sandbox.QiniuInjection{APIKey: &key}}
	if !hasAPIKey(spec) {
		t.Fatal("hasAPIKey() = false, want true for qiniu with api key")
	}

	spec = sandbox.InjectionSpec{HTTP: &sandbox.HTTPInjection{BaseURL: "https://api.example.com"}}
	if hasAPIKey(spec) {
		t.Fatal("hasAPIKey() = true, want false for http injection")
	}
}

func TestOptionalValue(t *testing.T) {
	if got := optionalValue(nil, "fallback"); got != "fallback" {
		t.Fatalf("optionalValue(nil) = %q, want %q", got, "fallback")
	}
	raw := "  value  "
	if got := optionalValue(&raw, "fallback"); got != "value" {
		t.Fatalf("optionalValue(trimmed) = %q, want %q", got, "value")
	}
}
