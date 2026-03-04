package operations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// LogsInfo holds parameters for viewing sandbox logs.
type LogsInfo struct {
	SandboxID string
	Level     string // Log level filter: INFO, WARN, ERROR, DEBUG
	Limit     int32
	Format    string // pretty or json
}

// Logs retrieves and displays sandbox logs.
func Logs(info LogsInfo) {
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

	params := &sandbox.GetLogsParams{}
	if info.Limit > 0 {
		params.Limit = &info.Limit
	}

	logs, err := sb.GetLogs(ctx, params)
	if err != nil {
		fmt.Printf("Error: get sandbox logs failed: %v\n", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(logs)
		return
	}

	if len(logs.LogEntries) > 0 {
		for _, entry := range logs.LogEntries {
			if info.Level != "" && string(entry.Level) != strings.ToUpper(info.Level) {
				continue
			}
			fmt.Printf("[%s] %s %s\n",
				entry.Timestamp.Format(time.RFC3339),
				entry.Level,
				entry.Message,
			)
		}
	} else if len(logs.Logs) > 0 {
		for _, l := range logs.Logs {
			fmt.Printf("[%s] %s\n",
				l.Timestamp.Format(time.RFC3339),
				l.Line,
			)
		}
	} else {
		fmt.Println("No logs found")
	}
}
