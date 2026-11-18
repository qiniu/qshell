package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// PauseInfo holds parameters for pausing sandboxes.
type PauseInfo struct {
	SandboxIDs []string
	// Memory 是否保存内存；nil 使用服务端默认值 true。
	Memory      *bool
	All         bool
	State       string // Comma-separated states for filtering when --all is used
	Metadata    string // Metadata filter: key=value
	MetadataSet bool   // Whether --metadata was explicitly supplied
}

// Pause pauses one or more sandboxes so they can be resumed later.
func Pause(info PauseInfo) {
	client, err := sbClient.NewSandboxListClient()
	if err != nil {
		data.SetCmdStatusError()
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	sandboxIDs := info.SandboxIDs

	// If --all flag is set, list and pause all matching sandboxes
	if info.All {
		params := &sandbox.ListParams{}
		// Default to "running" state when using --all (only running sandboxes can be paused)
		stateStr := info.State
		if stateStr == "" {
			stateStr = sbClient.DefaultState
		}
		states := sbClient.ParseStates(stateStr)
		params.State = &states

		if info.MetadataSet || info.Metadata != "" {
			filter, err := parseBulkMetadata(info.Metadata)
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
		fmt.Println("No sandboxes to pause")
		return
	}

	runSandboxBatch(sandboxIDs, "pause", "Paused", func(id string) error {
		return client.PauseWithParams(ctx, id, sandbox.PauseParams{Memory: info.Memory})
	})
}
