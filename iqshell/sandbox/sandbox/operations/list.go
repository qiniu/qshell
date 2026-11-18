package operations

import (
	"context"
	"fmt"
	"os"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/common/data"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ListInfo holds parameters for listing sandboxes.
type ListInfo struct {
	State       string // Comma-separated states: running,paused
	Metadata    string // Metadata filter: key=value
	MetadataSet bool   // Whether --metadata was explicitly supplied
	Limit       int32
	Format      string // pretty or json
	All         bool
	NextToken   string
}

// List lists sandboxes with optional filters.
func List(info ListInfo) {
	if info.All && info.NextToken != "" {
		data.SetCmdStatusError()
		sbClient.PrintError("--all and --next-token cannot be used together")
		return
	}
	client, err := sbClient.NewSandboxListClient()
	if err != nil {
		data.SetCmdStatusError()
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

	if info.MetadataSet || info.Metadata != "" {
		filter, err := parseMetadataFilter(info.Metadata)
		if err != nil {
			data.SetCmdStatusError()
			sbClient.PrintError("%v", err)
			return
		}
		params.Metadata = &filter
	}
	if info.Limit > 0 {
		params.Limit = &info.Limit
	}
	if info.NextToken != "" {
		params.NextToken = &info.NextToken
	}

	var sandboxes []sandbox.ListedSandbox
	var nextToken string
	if info.All {
		sandboxes, err = client.ListAll(context.Background(), params)
	} else {
		sandboxes, nextToken, err = client.ListPage(context.Background(), params)
	}
	if err != nil {
		data.SetCmdStatusError()
		sbClient.PrintError("list sandboxes failed: %v", err)
		return
	}
	if nextToken != "" {
		fmt.Fprintf(os.Stderr, "Next token: %s\n", nextToken)
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
