package operations

import (
	"context"
	"fmt"
	"strings"

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
		if info.State != "" {
			parts := strings.Split(info.State, ",")
			states := make([]sandbox.SandboxState, 0, len(parts))
			for _, s := range parts {
				s = strings.TrimSpace(s)
				if s != "" {
					states = append(states, sandbox.SandboxState(s))
				}
			}
			params.State = &states
		}
		if info.Metadata != "" {
			params.Metadata = &info.Metadata
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

	for _, id := range sandboxIDs {
		sb, cErr := client.Connect(ctx, id, sandbox.ConnectParams{Timeout: 10})
		if cErr != nil {
			fmt.Printf("Error: connect to sandbox %s failed: %v\n", id, cErr)
			continue
		}
		if kErr := sb.Kill(ctx); kErr != nil {
			fmt.Printf("Error: kill sandbox %s failed: %v\n", id, kErr)
			continue
		}
		fmt.Printf("Killed sandbox %s\n", id)
	}
}
