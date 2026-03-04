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
		fmt.Printf("Error: %v\n", err)
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
			fmt.Printf("Error: list sandboxes failed: %v\n", lErr)
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
				fmt.Printf("Error: connect to sandbox %s failed: %v\n", sandboxID, cErr)
				return
			}
			if kErr := sb.Kill(ctx); kErr != nil {
				fmt.Printf("Error: kill sandbox %s failed: %v\n", sandboxID, kErr)
				return
			}
			fmt.Printf("Killed sandbox %s\n", sandboxID)
		}(id)
	}
	wg.Wait()
}
