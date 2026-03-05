package operations

import (
	"context"
	"fmt"
	"os"

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
		sbClient.PrintError("%v", err)
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
		sbClient.PrintError("list sandboxes failed: %v", err)
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

	tw := sbClient.NewTable(os.Stdout)
	fmt.Fprintf(tw, "SANDBOX ID\tTEMPLATE ID\tALIAS\tSTARTED AT\tEND AT\tSTATE\tvCPUs\tRAM MiB\tENVD VERSION\tMETADATA\n")
	for _, sb := range sandboxes {
		var metadata string
		if sb.Metadata != nil {
			metadata = sbClient.FormatMetadata(map[string]string(*sb.Metadata))
		} else {
			metadata = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s\n",
			sb.SandboxID,
			sb.TemplateID,
			sbClient.FormatOptionalString(sb.Alias),
			sbClient.FormatTimestamp(sb.StartedAt),
			sbClient.FormatTimestamp(sb.EndAt),
			sb.State,
			sb.CPUCount,
			sb.MemoryMB,
			sb.EnvdVersion,
			metadata,
		)
	}
	tw.Flush()
}
