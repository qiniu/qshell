package operations

import (
	"context"
	"fmt"
	"sync"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ResumeInfo holds parameters for resuming sandboxes.
type ResumeInfo struct {
	SandboxIDs []string
	All        bool
	Metadata   string // Metadata filter: key=value
}

// Resume resumes one or more paused sandboxes.
// Uses the POST /sandboxes/{id}/resume API endpoint.
func Resume(info ResumeInfo) {
	// Still need the SDK client for --all listing
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	sandboxIDs := info.SandboxIDs

	// If --all flag is set, list and resume all paused sandboxes
	if info.All {
		params := &sandbox.ListParams{}
		// Default to "paused" state when using --all (only paused sandboxes need to be resumed)
		states := sbClient.ParseStates("paused")
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
		fmt.Println("No sandboxes to resume")
		return
	}

	// Resume sandboxes concurrently
	var wg sync.WaitGroup
	for _, id := range sandboxIDs {
		wg.Add(1)
		go func(sandboxID string) {
			defer wg.Done()
			if rErr := sbClient.ResumeSandbox(sandboxID, nil); rErr != nil {
				sbClient.PrintError("resume sandbox %s failed: %v", sandboxID, rErr)
				return
			}
			sbClient.PrintSuccess("Resumed sandbox %s", sandboxID)
		}(id)
	}
	wg.Wait()
}
