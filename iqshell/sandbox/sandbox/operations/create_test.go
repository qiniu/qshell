package operations

import "testing"

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
	injections, err := buildSandboxInjections(nil, []string{"type=http,base-url=https://api.example.com,headers=Authorization=Bearer token,X-Env=prod"})
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
