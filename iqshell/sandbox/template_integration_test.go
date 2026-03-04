//go:build integration

package sandbox

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

func TestIntegrationTemplateList(t *testing.T) {
	client := testSandboxClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	templates, err := client.ListTemplates(ctx, nil)
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	t.Logf("found %d template(s)", len(templates))
	for _, tmpl := range templates {
		t.Logf("  - %s (status=%s, aliases=%v)", tmpl.TemplateID, tmpl.BuildStatus, tmpl.Aliases)
	}
}

func TestIntegrationTemplateGet(t *testing.T) {
	client := testSandboxClient(t)
	templateID := findReadyTemplate(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tmpl, err := client.GetTemplate(ctx, templateID, nil)
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}
	if tmpl.TemplateID != templateID {
		t.Fatalf("TemplateID = %s, want %s", tmpl.TemplateID, templateID)
	}
	t.Logf("template %s: aliases=%v, public=%v, spawnCount=%d, builds=%d",
		tmpl.TemplateID, tmpl.Aliases, tmpl.Public, tmpl.SpawnCount, len(tmpl.Builds))

	if len(tmpl.Builds) > 0 {
		b := tmpl.Builds[0]
		t.Logf("  latest build: id=%s, status=%s, cpu=%d, memory=%dMB",
			b.BuildID, b.Status, b.CPUCount, b.MemoryMB)
	}
}

func TestIntegrationTemplateBuildStatus(t *testing.T) {
	client := testSandboxClient(t)
	templateID := findReadyTemplate(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get template to find a build ID
	tmpl, err := client.GetTemplate(ctx, templateID, nil)
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}
	if len(tmpl.Builds) == 0 {
		t.Skip("no builds found for template, skipping")
	}

	buildID := tmpl.Builds[0].BuildID
	t.Logf("checking build status for template=%s, build=%s", templateID, buildID)

	var buildInfo *sandbox.TemplateBuildInfo
	buildInfo, err = client.GetTemplateBuildStatus(ctx, templateID, buildID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "502") {
			t.Skipf("GetTemplateBuildStatus returned 502 (server-side issue), skipping")
		}
		t.Fatalf("GetTemplateBuildStatus failed: %v", err)
	}
	if buildInfo.TemplateID != templateID {
		t.Errorf("TemplateID = %s, want %s", buildInfo.TemplateID, templateID)
	}
	if buildInfo.BuildID != buildID {
		t.Errorf("BuildID = %s, want %s", buildInfo.BuildID, buildID)
	}
	if buildInfo.Status == "" {
		t.Error("Status should not be empty")
	}
	t.Logf("build status: %s (logs: %d lines)", buildInfo.Status, len(buildInfo.Logs))
}
