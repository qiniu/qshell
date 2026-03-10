package operations

import (
	"context"
	"fmt"
	"sync"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// PauseInfo holds parameters for pausing sandboxes.
type PauseInfo struct {
	SandboxIDs []string
	All        bool
	State      string // Comma-separated states for filtering when --all is used
	Metadata   string // Metadata filter: key=value
}

// Pause pauses one or more sandboxes so they can be resumed later.
func Pause(info PauseInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
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

		if info.Metadata != "" {
			m := sbClient.ParseMetadata(info.Metadata)
			if m != "" {
				params.Metadata = &m
			}
		}

		sandboxes, lErr := client.List(ctx, params)
		if lErr != nil {
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

	// Pause sandboxes concurrently
	var wg sync.WaitGroup
	for _, id := range sandboxIDs {
		wg.Add(1)
		go func(sandboxID string) {
			defer wg.Done()
			sb, cErr := client.Connect(ctx, sandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutCommand})
			if cErr != nil {
				sbClient.PrintError("connect to sandbox %s failed: %v", sandboxID, cErr)
				return
			}
			if pErr := sb.Pause(ctx); pErr != nil {
				sbClient.PrintError("pause sandbox %s failed: %v", sandboxID, pErr)
				return
			}
			sbClient.PrintSuccess("Paused sandbox %s", sandboxID)
		}(id)
	}
	wg.Wait()
}
