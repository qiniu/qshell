package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// MetricsInfo holds parameters for viewing sandbox metrics.
type MetricsInfo struct {
	SandboxID string
	Format    string // pretty or json
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

	ctx := context.Background()
	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutCommand})
	if err != nil {
		fmt.Printf("Error: connect to sandbox %s failed: %v\n", info.SandboxID, err)
		return
	}

	metrics, err := sb.GetMetrics(ctx, nil)
	if err != nil {
		fmt.Printf("Error: get sandbox metrics failed: %v\n", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(metrics)
		return
	}

	if len(metrics) == 0 {
		fmt.Println("No metrics available")
		return
	}

	fmt.Printf("%-6s %-10s %-15s %-15s %-15s %-15s %s\n",
		"CPU", "CPU %", "MEM USED", "MEM TOTAL", "DISK USED", "DISK TOTAL", "TIMESTAMP")
	for _, m := range metrics {
		fmt.Printf("%-6d %-10.1f %-15s %-15s %-15s %-15s %s\n",
			m.CPUCount,
			m.CPUUsedPct,
			formatBytes(m.MemUsed),
			formatBytes(m.MemTotal),
			formatBytes(m.DiskUsed),
			formatBytes(m.DiskTotal),
			m.Timestamp.Format(time.RFC3339),
		)
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
