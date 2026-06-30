//go:build integration

package sandbox

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

func testInjectionRuleClient(t *testing.T) *sandbox.Client {
	t.Helper()
	client, err := NewInjectionRuleClient()
	if err != nil {
		t.Skipf("injection rule client not available: %v", err)
	}
	return client
}

func TestIntegrationInjectionRuleMatchConditionsAndGithubBaseURL(t *testing.T) {
	client := testInjectionRuleClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	name := fmt.Sprintf("qshell-int-v72614-%d", time.Now().UnixNano())
	headers := map[string]string{"Authorization": "Bearer integration-test"}
	ifHeaders := map[string]string{"X-Qshell-Integration": "v7.26.14"}
	ifQueries := map[string]string{"inject": "true"}

	rule, err := client.CreateInjectionRule(ctx, sandbox.CreateInjectionRuleParams{
		Name: name,
		Injection: sandbox.InjectionSpec{
			HTTP: &sandbox.HTTPInjection{
				BaseURL:   "https://httpbin.org/v1/*",
				Headers:   &headers,
				IfHeaders: &ifHeaders,
				IfQueries: &ifQueries,
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateInjectionRule with HTTP match conditions failed: %v", err)
	}
	t.Logf("created injection rule %s", rule.RuleID)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if err := client.DeleteInjectionRule(cleanupCtx, rule.RuleID); err != nil {
			t.Logf("cleanup injection rule %s failed: %v", rule.RuleID, err)
		}
	})

	got, err := client.GetInjectionRule(ctx, rule.RuleID)
	if err != nil {
		t.Fatalf("GetInjectionRule after create failed: %v", err)
	}
	assertIntegrationHTTPInjection(t, got.Injection, headers, ifHeaders, ifQueries)

	githubToken := "ghp_qshell_integration_test"
	githubBaseURL := "https://api.github.com/repos/qiniu/*"
	githubIfHeaders := map[string]string{"X-GitHub-Api-Version": "2022-11-28"}
	githubIfQueries := map[string]string{"per_page": "100"}

	updated, err := client.UpdateInjectionRule(ctx, rule.RuleID, sandbox.UpdateInjectionRuleParams{
		Injection: &sandbox.InjectionSpec{
			Github: &sandbox.GithubInjection{
				BaseURL:   &githubBaseURL,
				IfHeaders: &githubIfHeaders,
				IfQueries: &githubIfQueries,
				Token:     &githubToken,
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateInjectionRule with GitHub base URL and match conditions failed: %v", err)
	}
	assertIntegrationGithubInjection(t, updated.Injection, githubBaseURL, githubIfHeaders, githubIfQueries)

	got, err = client.GetInjectionRule(ctx, rule.RuleID)
	if err != nil {
		t.Fatalf("GetInjectionRule after update failed: %v", err)
	}
	assertIntegrationGithubInjection(t, got.Injection, githubBaseURL, githubIfHeaders, githubIfQueries)

	rules, err := client.ListInjectionRules(ctx)
	if err != nil {
		t.Fatalf("ListInjectionRules failed: %v", err)
	}
	found := false
	for _, candidate := range rules {
		if candidate.RuleID == rule.RuleID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created injection rule %s not found in list response", rule.RuleID)
	}
}

func assertIntegrationHTTPInjection(t *testing.T, spec sandbox.InjectionSpec, headers, ifHeaders, ifQueries map[string]string) {
	t.Helper()
	if spec.HTTP == nil {
		t.Fatalf("HTTP injection is nil: %+v", spec)
	}
	if spec.HTTP.BaseURL != "https://httpbin.org/v1/*" {
		t.Fatalf("HTTP base URL = %q, want https://httpbin.org/v1/*", spec.HTTP.BaseURL)
	}
	assertStringMapPtrEqual(t, "HTTP headers", spec.HTTP.Headers, headers)
	assertStringMapPtrEqual(t, "HTTP if headers", spec.HTTP.IfHeaders, ifHeaders)
	assertStringMapPtrEqual(t, "HTTP if queries", spec.HTTP.IfQueries, ifQueries)
}

func assertIntegrationGithubInjection(t *testing.T, spec sandbox.InjectionSpec, baseURL string, ifHeaders, ifQueries map[string]string) {
	t.Helper()
	if spec.Github == nil {
		t.Fatalf("GitHub injection is nil: %+v", spec)
	}
	if spec.Github.BaseURL == nil || *spec.Github.BaseURL != baseURL {
		t.Fatalf("GitHub base URL = %v, want %q", spec.Github.BaseURL, baseURL)
	}
	assertStringMapPtrEqual(t, "GitHub if headers", spec.Github.IfHeaders, ifHeaders)
	assertStringMapPtrEqual(t, "GitHub if queries", spec.Github.IfQueries, ifQueries)
	if spec.Github.Token == nil || *spec.Github.Token == "" {
		t.Fatal("GitHub token was not returned as configured")
	}
}

func assertStringMapPtrEqual(t *testing.T, label string, got *map[string]string, want map[string]string) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s = nil, want %v", label, want)
	}
	if len(*got) != len(want) {
		t.Fatalf("%s len = %d, want %d: %v", label, len(*got), len(want), *got)
	}
	for key, wantValue := range want {
		if gotValue := (*got)[key]; gotValue != wantValue {
			t.Fatalf("%s[%q] = %q, want %q", label, key, gotValue, wantValue)
		}
	}
}
