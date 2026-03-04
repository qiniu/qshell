package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		fmt.Println("Error: sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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
		fmt.Printf("Error: connect to sandbox %s failed: %v\n", info.SandboxID, err)
		return
	}

	var lastTimestamp *time.Time
	headerPrinted := false

	for {
		params := &sandbox.GetMetricsParams{}
		if lastTimestamp != nil {
			start := lastTimestamp.Unix()
			params.Start = &start
		}

		metrics, mErr := sb.GetMetrics(ctx, params)
		if mErr != nil {
			fmt.Printf("Error: get sandbox metrics failed: %v\n", mErr)
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
				if !headerPrinted {
					fmt.Printf("%-6s %-10s %-15s %-15s %-15s %-15s %s\n",
						"CPU", "CPU %", "MEM USED", "MEM TOTAL", "DISK USED", "DISK TOTAL", "TIMESTAMP")
					headerPrinted = true
				}
				fmt.Printf("%-6d %-10.1f %-15s %-15s %-15s %-15s %s\n",
					m.CPUCount,
					m.CPUUsedPct,
					formatBytes(m.MemUsed),
					formatBytes(m.MemTotal),
					formatBytes(m.DiskUsed),
					formatBytes(m.DiskTotal),
					m.Timestamp.Format(time.RFC3339),
				)
				ts := m.Timestamp
				lastTimestamp = &ts
			}
		}

		if !info.Follow {
			return
		}

		// Check if sandbox is still running
		running, rErr := sb.IsRunning(ctx)
		if rErr != nil || !running {
			if info.Format != sbClient.FormatJSON {
				fmt.Println("\nStopped printing metrics — sandbox is closed")
			}
			return
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(400 * time.Millisecond)
	}
}

// formatBytes formats byte count to human-readable string.
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
