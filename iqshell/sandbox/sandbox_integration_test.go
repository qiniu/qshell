//go:build integration

package sandbox

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

// testSandboxClient creates a client from environment variables; skips if no API key is set.
func testSandboxClient(t *testing.T) *sandbox.Client {
	t.Helper()
	client, err := NewSandboxClient()
	if err != nil {
		t.Skipf("sandbox client not available: %v", err)
	}
	return client
}

// findReadyTemplate returns the first ready/uploaded template ID; skips if none found.
func findReadyTemplate(t *testing.T, client *sandbox.Client) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	templates, err := client.ListTemplates(ctx, nil)
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	for _, tmpl := range templates {
		if tmpl.BuildStatus == sandbox.BuildStatusReady || tmpl.BuildStatus == sandbox.BuildStatusUploaded {
			return tmpl.TemplateID
		}
	}
	t.Skip("no ready template available, skipping")
	return ""
}

// killAllRunning kills all currently running sandboxes to free up quota.
func killAllRunning(t *testing.T, client *sandbox.Client) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	states := []sandbox.SandboxState{sandbox.StateRunning}
	sandboxes, err := client.List(ctx, &sandbox.ListParams{State: &states})
	if err != nil {
		t.Fatalf("List running sandboxes failed: %v", err)
	}
	for _, s := range sandboxes {
		sb, cErr := client.Connect(ctx, s.SandboxID, sandbox.ConnectParams{Timeout: ConnectTimeoutCommand})
		if cErr != nil {
			t.Logf("connect to %s for cleanup failed: %v", s.SandboxID, cErr)
			continue
		}
		if kErr := sb.Kill(ctx); kErr != nil {
			t.Logf("kill %s for cleanup failed: %v", s.SandboxID, kErr)
		}
	}
	if len(sandboxes) > 0 {
		t.Logf("cleaned up %d running sandbox(es) to free quota", len(sandboxes))
	}
}

// createTestSandbox creates a sandbox and registers t.Cleanup to kill it automatically.
func createTestSandbox(t *testing.T, client *sandbox.Client) *sandbox.Sandbox {
	t.Helper()
	templateID := findReadyTemplate(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	timeout := int32(120)
	sb, _, err := client.CreateAndWait(ctx, sandbox.CreateParams{
		TemplateID: templateID,
		Timeout:    &timeout,
	}, sandbox.WithPollInterval(2*time.Second))
	if err != nil {
		t.Fatalf("CreateAndWait failed: %v", err)
	}
	t.Logf("created sandbox %s (template=%s)", sb.ID(), templateID)

	t.Cleanup(func() {
		killCtx, killCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer killCancel()
		if err := sb.Kill(killCtx); err != nil {
			t.Logf("cleanup sandbox %s failed: %v", sb.ID(), err)
		}
	})
	return sb
}

// --- Tests that do NOT create sandboxes (always safe to run) ---

func TestIntegrationSandboxList(t *testing.T) {
	client := testSandboxClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	states := ParseStates(DefaultState)
	sandboxes, err := client.List(ctx, &sandbox.ListParams{
		State: &states,
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	t.Logf("found %d running sandbox(es)", len(sandboxes))
}

func TestIntegrationSandboxListJSON(t *testing.T) {
	client := testSandboxClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sandboxes, err := client.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	data, err := json.MarshalIndent(sandboxes, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	if !json.Valid(data) {
		t.Fatal("List result is not valid JSON")
	}
	t.Logf("JSON output: %d bytes, %d sandbox(es)", len(data), len(sandboxes))
}

// --- Tests that share ONE sandbox (Connect, Logs, Metrics, SetTimeout, ListWithMetadata) ---

func TestIntegrationSandboxShared(t *testing.T) {
	client := testSandboxClient(t)

	// Clean up existing sandboxes to free quota before creating ours
	killAllRunning(t, client)

	templateID := findReadyTemplate(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Create a single shared sandbox with metadata
	timeout := int32(120)
	meta := sandbox.Metadata{"inttest": "shared"}
	sb, _, err := client.CreateAndWait(ctx, sandbox.CreateParams{
		TemplateID: templateID,
		Timeout:    &timeout,
		Metadata:   &meta,
	}, sandbox.WithPollInterval(2*time.Second))
	if err != nil {
		t.Fatalf("CreateAndWait failed: %v", err)
	}
	t.Logf("shared sandbox created: %s", sb.ID())
	t.Cleanup(func() {
		killCtx, killCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer killCancel()
		if err := sb.Kill(killCtx); err != nil {
			t.Logf("cleanup shared sandbox %s failed: %v", sb.ID(), err)
		}
	})

	t.Run("Connect", func(t *testing.T) {
		connectCtx, connectCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer connectCancel()

		connected, err := client.Connect(connectCtx, sb.ID(), sandbox.ConnectParams{Timeout: ConnectTimeoutCommand})
		if err != nil {
			t.Fatalf("Connect failed: %v", err)
		}
		if connected.ID() != sb.ID() {
			t.Fatalf("Connect returned ID %s, want %s", connected.ID(), sb.ID())
		}
		t.Logf("connected to sandbox %s", connected.ID())
	})

	t.Run("Logs", func(t *testing.T) {
		logsCtx, logsCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer logsCancel()

		logs, err := sb.GetLogs(logsCtx, &sandbox.GetLogsParams{})
		if err != nil {
			if strings.Contains(err.Error(), "502") {
				t.Skipf("GetLogs returned 502 (server-side issue), skipping")
			}
			t.Fatalf("GetLogs failed: %v", err)
		}
		t.Logf("got %d logs, %d log entries", len(logs.Logs), len(logs.LogEntries))

		for _, entry := range logs.LogEntries {
			level := string(entry.Level)
			if !IsLogLevelIncluded(level, "DEBUG") {
				t.Errorf("log entry level %q should be included at DEBUG minimum", level)
			}
		}
	})

	t.Run("Metrics", func(t *testing.T) {
		// Wait briefly for metrics to be available
		time.Sleep(2 * time.Second)

		metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer metricsCancel()

		metrics, err := sb.GetMetrics(metricsCtx, &sandbox.GetMetricsParams{})
		if err != nil {
			t.Fatalf("GetMetrics failed: %v", err)
		}
		t.Logf("got %d metric(s)", len(metrics))

		if len(metrics) > 0 {
			m := metrics[len(metrics)-1]
			if m.CPUCount <= 0 {
				t.Errorf("CPUCount = %d, want > 0", m.CPUCount)
			}
			if m.MemTotal <= 0 {
				t.Errorf("MemTotal = %d, want > 0", m.MemTotal)
			}
			t.Logf("latest metric: cpu=%d, cpuPct=%.1f%%, memUsed=%d, memTotal=%d",
				m.CPUCount, m.CPUUsedPct, m.MemUsed, m.MemTotal)
		}
	})

	t.Run("SetTimeout", func(t *testing.T) {
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer timeoutCancel()

		if err := sb.SetTimeout(timeoutCtx, 30*time.Second); err != nil {
			t.Fatalf("SetTimeout(30s) failed: %v", err)
		}
		t.Log("SetTimeout(30s) succeeded")
	})

	t.Run("ListWithMetadata", func(t *testing.T) {
		listCtx, listCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer listCancel()

		metadataQuery := ParseMetadata("inttest=shared")
		if metadataQuery == "" {
			t.Fatal("ParseMetadata returned empty string")
		}

		sandboxes, err := client.List(listCtx, &sandbox.ListParams{
			Metadata: &metadataQuery,
		})
		if err != nil {
			t.Fatalf("List with metadata failed: %v", err)
		}

		found := false
		for _, s := range sandboxes {
			if s.SandboxID == sb.ID() {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("sandbox %s not found in metadata-filtered list (%d results)", sb.ID(), len(sandboxes))
		}
		t.Logf("sandbox %s found in metadata-filtered list", sb.ID())
	})
}

// --- Create-and-Kill lifecycle test (creates 1 sandbox, kills it immediately) ---

func TestIntegrationSandboxCreateAndKill(t *testing.T) {
	client := testSandboxClient(t)

	// Clean up to ensure quota is available
	killAllRunning(t, client)

	templateID := findReadyTemplate(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Create
	timeout := int32(120)
	sb, _, err := client.CreateAndWait(ctx, sandbox.CreateParams{
		TemplateID: templateID,
		Timeout:    &timeout,
	}, sandbox.WithPollInterval(2*time.Second))
	if err != nil {
		t.Fatalf("CreateAndWait failed: %v", err)
	}
	t.Logf("created sandbox %s", sb.ID())

	// Confirm in list
	sandboxes, err := client.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	found := false
	for _, s := range sandboxes {
		if s.SandboxID == sb.ID() {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("sandbox not found in list after creation")
	}

	// Kill
	if err := sb.Kill(ctx); err != nil {
		t.Fatalf("Kill failed: %v", err)
	}
	t.Logf("killed sandbox %s", sb.ID())

	// Confirm not in running list
	states := []sandbox.SandboxState{sandbox.StateRunning}
	sandboxes, err = client.List(ctx, &sandbox.ListParams{
		State: &states,
	})
	if err != nil {
		t.Fatalf("List after kill failed: %v", err)
	}
	for _, s := range sandboxes {
		if s.SandboxID == sb.ID() {
			t.Fatal("sandbox still in running list after kill")
		}
	}
	t.Log("sandbox confirmed removed from running list")
}

// --- KillAll test (creates 2 sandboxes, kills them concurrently) ---

func TestIntegrationSandboxKillAll(t *testing.T) {
	client := testSandboxClient(t)

	// Clean up all existing sandboxes first to ensure we have quota for 2
	killAllRunning(t, client)

	templateID := findReadyTemplate(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	timeout := int32(120)
	meta := sandbox.Metadata{"inttest": "killall"}

	sb1, _, err := client.CreateAndWait(ctx, sandbox.CreateParams{
		TemplateID: templateID,
		Timeout:    &timeout,
		Metadata:   &meta,
	}, sandbox.WithPollInterval(2*time.Second))
	if err != nil {
		t.Fatalf("create sandbox 1 failed: %v", err)
	}
	t.Logf("created sandbox 1: %s", sb1.ID())

	sb2, _, err := client.CreateAndWait(ctx, sandbox.CreateParams{
		TemplateID: templateID,
		Timeout:    &timeout,
		Metadata:   &meta,
	}, sandbox.WithPollInterval(2*time.Second))
	if err != nil {
		sb1.Kill(ctx)
		t.Fatalf("create sandbox 2 failed: %v", err)
	}
	t.Logf("created sandbox 2: %s", sb2.ID())

	// List with metadata filter
	metadataQuery := ParseMetadata("inttest=killall")
	sandboxes, err := client.List(ctx, &sandbox.ListParams{
		Metadata: &metadataQuery,
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(sandboxes) < 2 {
		t.Fatalf("expected at least 2 sandboxes, got %d", len(sandboxes))
	}

	// Kill all concurrently (same pattern as kill --all)
	var wg sync.WaitGroup
	for _, s := range sandboxes {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			connected, cErr := client.Connect(ctx, id, sandbox.ConnectParams{Timeout: ConnectTimeoutCommand})
			if cErr != nil {
				t.Errorf("connect to %s failed: %v", id, cErr)
				return
			}
			if kErr := connected.Kill(ctx); kErr != nil {
				t.Errorf("kill %s failed: %v", id, kErr)
				return
			}
			t.Logf("killed %s", id)
		}(s.SandboxID)
	}
	wg.Wait()

	// Verify all terminated
	states := []sandbox.SandboxState{sandbox.StateRunning}
	remaining, err := client.List(ctx, &sandbox.ListParams{
		State:    &states,
		Metadata: &metadataQuery,
	})
	if err != nil {
		t.Fatalf("List after kill-all failed: %v", err)
	}
	if len(remaining) > 0 {
		t.Fatalf("expected 0 running sandboxes after kill-all, got %d", len(remaining))
	}
	t.Log("all sandboxes killed successfully")
}
