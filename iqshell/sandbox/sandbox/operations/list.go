package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ListInfo holds parameters for listing sandboxes.
type ListInfo struct {
	State    string // Comma-separated states: running,paused
	Metadata string // Metadata filter: key=value
	Limit    int32
	Format   string // pretty or json
}

// List lists sandboxes with optional filters.
func List(info ListInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	params := &sandbox.ListParams{}
	// Default to "running" state if not specified (matches e2b CLI behavior)
	stateStr := info.State
	if stateStr == "" {
		stateStr = sbClient.DefaultState
	}
	states := sbClient.ParseStates(stateStr)
	params.State = &states

	if info.Metadata != "" {
		m := sbClient.ParseMetadata(info.Metadata)
		if m != "" {
			params.Metadata = &m
		}
	}
	if info.Limit > 0 {
		params.Limit = &info.Limit
	}

	sandboxes, err := client.List(context.Background(), params)
	if err != nil {
		fmt.Printf("Error: list sandboxes failed: %v\n", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(sandboxes)
		return
	}

	if len(sandboxes) == 0 {
		fmt.Println("No sandboxes found")
		return
	}

	fmt.Printf("%-30s %-20s %-10s %-6s %-10s %-10s %-22s %s\n",
		"SANDBOX ID", "TEMPLATE ID", "STATE", "CPU", "MEMORY", "DISK", "STARTED AT", "END AT")
	for _, sb := range sandboxes {
		fmt.Printf("%-30s %-20s %-10s %-6d %-10s %-10s %-22s %s\n",
			sb.SandboxID,
			sb.TemplateID,
			sb.State,
			sb.CPUCount,
			fmt.Sprintf("%dMB", sb.MemoryMB),
			fmt.Sprintf("%dMB", sb.DiskSizeMB),
			sb.StartedAt.Format(time.RFC3339),
			sb.EndAt.Format(time.RFC3339),
		)
	}
}
