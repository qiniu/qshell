package operations

import (
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

func TestBuildSandboxInjections_Empty(t *testing.T) {
	injections, err := buildSandboxInjections(nil, nil)
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if injections != nil {
		t.Fatalf("buildSandboxInjections() = %v, want nil", injections)
	}
}

func TestBuildSandboxInjections_WithRuleIDs(t *testing.T) {
	injections, err := buildSandboxInjections([]string{"rule-1", " rule-2 "}, nil)
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if len(injections) != 2 {
		t.Fatalf("buildSandboxInjections() len = %d, want 2", len(injections))
	}
	if injections[0].ByID == nil || *injections[0].ByID != "rule-1" {
		t.Fatalf("first injection = %+v, want rule-1", injections[0])
	}
	if injections[1].ByID == nil || *injections[1].ByID != "rule-2" {
		t.Fatalf("second injection = %+v, want rule-2", injections[1])
	}
}

func TestBuildSandboxInjections_RejectsEmptyRuleID(t *testing.T) {
	if _, err := buildSandboxInjections([]string{"rule-1", " "}, nil); err == nil {
		t.Fatal("expected empty injection rule ID to fail")
	}
}

func TestBuildSandboxInjections_WithInlineOpenAI(t *testing.T) {
	injections, err := buildSandboxInjections(nil, []string{"type=openai,api-key=sk-test,base-url=https://api.openai-proxy.example.com"})
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if len(injections) != 1 {
		t.Fatalf("buildSandboxInjections() len = %d, want 1", len(injections))
	}
	if injections[0].OpenAI == nil {
		t.Fatalf("first injection = %+v, want openai injection", injections[0])
	}
	if injections[0].OpenAI.APIKey == nil || *injections[0].OpenAI.APIKey != "sk-test" {
		t.Fatalf("openai api key = %v, want sk-test", injections[0].OpenAI.APIKey)
	}
}

func TestBuildSandboxInjections_WithInlineHTTP(t *testing.T) {
	injections, err := buildSandboxInjections(nil, []string{"type=http,base-url=https://api.example.com,headers=Authorization=Bearer token;X-Env=prod"})
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if len(injections) != 1 {
		t.Fatalf("buildSandboxInjections() len = %d, want 1", len(injections))
	}
	if injections[0].HTTP == nil {
		t.Fatalf("first injection = %+v, want http injection", injections[0])
	}
	if injections[0].HTTP.Headers == nil || len(*injections[0].HTTP.Headers) != 2 {
		t.Fatalf("http headers = %v, want 2 headers", injections[0].HTTP.Headers)
	}
}

func TestBuildSandboxInjections_WithRulesAndInline(t *testing.T) {
	injections, err := buildSandboxInjections(
		[]string{"rule-1"},
		[]string{"type=gemini,api-key=sk-gem"},
	)
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if len(injections) != 2 {
		t.Fatalf("buildSandboxInjections() len = %d, want 2", len(injections))
	}
	if injections[0].ByID == nil || *injections[0].ByID != "rule-1" {
		t.Fatalf("first injection = %+v, want by-id rule-1", injections[0])
	}
	if injections[1].Gemini == nil {
		t.Fatalf("second injection = %+v, want gemini injection", injections[1])
	}
}

func TestBuildSandboxInjections_RejectsInvalidInlineSpec(t *testing.T) {
	if _, err := buildSandboxInjections(nil, []string{"api-key=sk-test"}); err == nil {
		t.Fatal("expected inline injection without type to fail")
	}
}

func TestBuildSandboxInjections_RejectsInvalidInlineHTTPURL(t *testing.T) {
	if _, err := buildSandboxInjections(nil, []string{"type=http,base-url=not-a-url"}); err == nil {
		t.Fatal("expected inline http injection with invalid URL to fail")
	}
}

func TestBuildSandboxInjections_RejectsUnsupportedInlineType(t *testing.T) {
	if _, err := buildSandboxInjections(nil, []string{"type=unknown"}); err == nil {
		t.Fatal("expected unsupported inline injection type to fail")
	}
}

func TestBuildSandboxInjections_WithInlineQiniu(t *testing.T) {
	injections, err := buildSandboxInjections(nil, []string{"type=qiniu,api-key=sk-qiniu,base-url=https://api.qnaigc-proxy.example.com"})
	if err != nil {
		t.Fatalf("buildSandboxInjections() error = %v", err)
	}
	if len(injections) != 1 {
		t.Fatalf("buildSandboxInjections() len = %d, want 1", len(injections))
	}
	if injections[0].Qiniu == nil {
		t.Fatalf("first injection = %+v, want qiniu injection", injections[0])
	}
	if injections[0].Qiniu.APIKey == nil || *injections[0].Qiniu.APIKey != "sk-qiniu" {
		t.Fatalf("qiniu api key = %v, want sk-qiniu", injections[0].Qiniu.APIKey)
	}
}

func TestParseInlineInjectionFields_HeadersOnly(t *testing.T) {
	fields := parseInlineInjectionFields("headers=Authorization=Bearer token;X-Env=prod")
	if got := fields["headers"]; got != "Authorization=Bearer token;X-Env=prod" {
		t.Fatalf("headers = %q, want %q", got, "Authorization=Bearer token;X-Env=prod")
	}
}

func TestParseInlineInjectionFields_HeadersWithOtherFields(t *testing.T) {
	fields := parseInlineInjectionFields("type=http,base-url=https://api.example.com,headers=Authorization=Bearer token;X-Env=prod")
	if fields["type"] != "http" || fields["base-url"] != "https://api.example.com" {
		t.Fatalf("fields = %v, want type and base-url parsed", fields)
	}
	if got := fields["headers"]; got != "Authorization=Bearer token;X-Env=prod" {
		t.Fatalf("headers = %q, want %q", got, "Authorization=Bearer token;X-Env=prod")
	}
}

func TestParseInlineHeaders_SemicolonSeparated(t *testing.T) {
	headers := parseInlineHeaders("Authorization=Bearer token;X-Env=prod")
	if len(headers) != 2 {
		t.Fatalf("headers = %v, want 2 headers", headers)
	}
	if headers["Authorization"] != "Bearer token" || headers["X-Env"] != "prod" {
		t.Fatalf("headers = %v, want parsed headers", headers)
	}
}

func TestParseInlineHeaders_CommaFallback(t *testing.T) {
	headers := parseInlineHeaders("Authorization=Bearer token,X-Env=prod")
	if len(headers) != 2 {
		t.Fatalf("headers = %v, want 2 headers", headers)
	}
	if headers["Authorization"] != "Bearer token" || headers["X-Env"] != "prod" {
		t.Fatalf("headers = %v, want parsed headers", headers)
	}
}

// === buildSandboxResources tests ===

func TestBuildSandboxResources_Empty(t *testing.T) {
	resources, err := buildSandboxResources(nil)
	if err != nil {
		t.Fatalf("buildSandboxResources() error = %v", err)
	}
	if resources != nil {
		t.Fatalf("buildSandboxResources() = %v, want nil", resources)
	}
}

func TestBuildSandboxResources_GithubRepository(t *testing.T) {
	resources, err := buildSandboxResources([]string{
		"type=github_repository,url=https://github.com/owner/repo.git,mount-path=/workspace/repo,token=ghp-xxx",
	})
	if err != nil {
		t.Fatalf("buildSandboxResources() error = %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("buildSandboxResources() len = %d, want 1", len(resources))
	}
	got := resources[0].GitRepository
	if got == nil {
		t.Fatalf("resource = %+v, want GitRepository set", resources[0])
	}
	if got.Type != sandbox.GitRepositoryTypeGithub {
		t.Fatalf("type = %q, want %q", got.Type, sandbox.GitRepositoryTypeGithub)
	}
	if got.URL != "https://github.com/owner/repo.git" {
		t.Fatalf("url = %q, want %q", got.URL, "https://github.com/owner/repo.git")
	}
	if got.MountPath != "/workspace/repo" {
		t.Fatalf("mount path = %q, want %q", got.MountPath, "/workspace/repo")
	}
	if got.AuthorizationToken == nil || *got.AuthorizationToken != "ghp-xxx" {
		t.Fatalf("token = %v, want ghp-xxx", got.AuthorizationToken)
	}
}

func TestBuildSandboxResources_DefaultsTypeAndAcceptsMountAlias(t *testing.T) {
	resources, err := buildSandboxResources([]string{
		"url=https://github.com/owner/repo.git,mount=/workspace/repo",
	})
	if err != nil {
		t.Fatalf("buildSandboxResources() error = %v", err)
	}
	got := resources[0].GitRepository
	if got == nil {
		t.Fatal("resource GitRepository = nil, want set when type omitted")
	}
	if got.Type != sandbox.GitRepositoryTypeGithub {
		t.Fatalf("type defaulted = %q, want %q", got.Type, sandbox.GitRepositoryTypeGithub)
	}
	if got.MountPath != "/workspace/repo" {
		t.Fatalf("mount path via mount= alias = %q, want /workspace/repo", got.MountPath)
	}
	if got.AuthorizationToken != nil {
		t.Fatalf("token = %v, want nil when not provided", got.AuthorizationToken)
	}
}

func TestBuildSandboxResources_RejectsMissingURL(t *testing.T) {
	if _, err := buildSandboxResources([]string{"type=github_repository,mount-path=/workspace"}); err == nil {
		t.Fatal("expected missing url to fail")
	}
}

func TestBuildSandboxResources_RejectsMissingMountPath(t *testing.T) {
	if _, err := buildSandboxResources([]string{"type=github_repository,url=https://github.com/owner/repo.git"}); err == nil {
		t.Fatal("expected missing mount-path to fail")
	}
}

func TestBuildSandboxResources_RejectsUnsupportedType(t *testing.T) {
	if _, err := buildSandboxResources([]string{"type=gitlab_repository,url=https://gitlab.com/owner/repo.git,mount-path=/workspace"}); err == nil {
		t.Fatal("expected unsupported resource type to fail")
	}
}

func TestBuildSandboxResources_Multiple(t *testing.T) {
	resources, err := buildSandboxResources([]string{
		"url=https://github.com/owner/a.git,mount-path=/workspace/a",
		"url=https://github.com/owner/b.git,mount-path=/workspace/b",
	})
	if err != nil {
		t.Fatalf("buildSandboxResources() error = %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("buildSandboxResources() len = %d, want 2", len(resources))
	}
}
