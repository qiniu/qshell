package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// MetricsInfo holds parameters for viewing sandbox metrics.
type MetricsInfo struct {
	SandboxID string
	Format    string // pretty or json
	Follow    bool   // Keep streaming metrics until sandbox is closed
}

// Metrics retrieves and displays sandbox resource metrics.
func Metrics(info MetricsInfo) {
	if info.SandboxID == "" {
		sbClient.PrintError("sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
	defer signal.Stop(sigCh)

	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutCommand})
	if err != nil {
		sbClient.PrintError("connect to sandbox %s failed: %v", info.SandboxID, err)
		return
	}

	// Async sandbox-done monitoring (non-blocking, checks every 5s)
	sandboxDone := make(chan struct{})
	if info.Follow {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					running, _ := sb.IsRunning(ctx)
					if !running {
						close(sandboxDone)
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	cyanLabel := color.New(color.FgCyan)
	var lastTimestamp *time.Time

	for {
		params := &sandbox.GetMetricsParams{}
		if lastTimestamp != nil {
			start := lastTimestamp.Unix()
			params.Start = &start
		}

		metrics, mErr := sb.GetMetrics(ctx, params)
		if mErr != nil {
			sbClient.PrintError("get sandbox metrics failed: %v", mErr)
			return
		}

		if info.Format == sbClient.FormatJSON {
			if !info.Follow {
				sbClient.PrintJSON(metrics)
				return
			}
			// In follow+json mode, print each batch
			if len(metrics) > 0 {
				sbClient.PrintJSON(metrics)
			}
		} else {
			if !info.Follow && len(metrics) == 0 {
				fmt.Println("No metrics available")
				return
			}
			for _, m := range metrics {
				// Skip metrics with same or earlier timestamp
				if lastTimestamp != nil && !m.Timestamp.After(*lastTimestamp) {
					continue
				}
				// e2b inline format:
				// [timestamp]  CPU:  12.5% / 2 Cores | Memory:  256 / 512 MiB | Disk:  100 / 1024 MiB
				fmt.Printf("[%s]  %s  %.1f%% / %d Cores | %s  %s / %s | %s  %s / %s\n",
					m.Timestamp.Format(time.RFC3339),
					cyanLabel.Sprint("CPU:"),
					m.CPUUsedPct,
					m.CPUCount,
					cyanLabel.Sprint("Memory:"),
					sbClient.FormatBytes(m.MemUsed),
					sbClient.FormatBytes(m.MemTotal),
					cyanLabel.Sprint("Disk:"),
					sbClient.FormatBytes(m.DiskUsed),
					sbClient.FormatBytes(m.DiskTotal),
				)
				ts := m.Timestamp
				lastTimestamp = &ts
			}
		}

		if !info.Follow {
			return
		}

		// Check if sandbox is done or context cancelled
		select {
		case <-sandboxDone:
			if info.Format != sbClient.FormatJSON {
				fmt.Println("\nStopped printing metrics — sandbox is closed")
			}
			return
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(400 * time.Millisecond)
	}
}
