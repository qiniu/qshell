package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/qiniu/qshell/v2/iqshell/common/data"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// KillInfo holds parameters for killing sandboxes.
type KillInfo struct {
	SandboxIDs  []string
	All         bool
	State       string // Comma-separated states for filtering when --all is used
	Metadata    string // Metadata filter: key=value
	MetadataSet bool   // Whether --metadata was explicitly supplied
}

// Kill terminates one or more sandboxes.
func Kill(info KillInfo) {
	client, err := sbClient.NewSandboxListClient()
	if err != nil {
		data.SetCmdStatusError()
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	sandboxIDs := info.SandboxIDs

	// If --all flag is set, list and kill all matching sandboxes
	if info.All {
		params := &sandbox.ListParams{}
		// Default to "running" state when using --all (matches e2b CLI behavior)
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

		sandboxes, lErr := client.ListAll(ctx, params)
		if lErr != nil {
			data.SetCmdStatusError()
			sbClient.PrintError("list sandboxes failed: %v", lErr)
			return
		}

		sandboxIDs = make([]string, 0, len(sandboxes))
		for _, sb := range sandboxes {
			sandboxIDs = append(sandboxIDs, sb.SandboxID)
		}
	}

	if len(sandboxIDs) == 0 {
		fmt.Println("No sandboxes to kill")
		return
	}

	runSandboxBatch(sandboxIDs, "kill", "Killed", func(id string) error {
		return client.Kill(ctx, id)
	})
}
