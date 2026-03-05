package operations

import (
	"context"
	"fmt"
	"sync"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// KillInfo holds parameters for killing sandboxes.
type KillInfo struct {
	SandboxIDs []string
	All        bool
	State      string // Comma-separated states for filtering when --all is used
	Metadata   string // Metadata filter: key=value
}

// Kill terminates one or more sandboxes.
func Kill(info KillInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
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
		fmt.Println("No sandboxes to kill")
		return
	}

	// Kill sandboxes concurrently
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
			if kErr := sb.Kill(ctx); kErr != nil {
				sbClient.PrintError("kill sandbox %s failed: %v", sandboxID, kErr)
				return
			}
			sbClient.PrintSuccess("Killed sandbox %s", sandboxID)
		}(id)
	}
	wg.Wait()
}
