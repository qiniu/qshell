package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ResumeInfo holds parameters for resuming sandboxes.
type ResumeInfo struct {
	SandboxIDs  []string
	All         bool
	Metadata    string // Metadata filter: key=value
	MetadataSet bool   // Whether --metadata was explicitly supplied
}

// Resume resumes one or more paused sandboxes.
// 使用 SDK v2 Connect 恢复沙箱，省略超时时间时采用服务端默认值。
func Resume(info ResumeInfo) {
	client, err := sbClient.NewSandboxListClient()
	if err != nil {
		data.SetCmdStatusError()
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
		fmt.Println("No sandboxes to resume")
		return
	}

	runSandboxBatch(sandboxIDs, "resume", "Resumed", func(id string) error {
		_, err := client.Connect(ctx, id, sandbox.ConnectParams{})
		return err
	})
}
