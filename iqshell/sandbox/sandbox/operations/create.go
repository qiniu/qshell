package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// CreateInfo holds parameters for creating a sandbox.
type CreateInfo struct {
	TemplateID string
	Timeout    int32
	Metadata   string
}

// Create creates a new sandbox and connects to its terminal.
// When the terminal session ends, the sandbox is killed.
// The sandbox stays alive via keep-alive in the terminal session (matches e2b CLI behavior).
func Create(info CreateInfo) {
	if info.TemplateID == "" {
		sbClient.PrintError("template ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	params := sandbox.CreateParams{
		TemplateID: info.TemplateID,
	}
	if info.Timeout > 0 {
		params.Timeout = &info.Timeout
	}
	if info.Metadata != "" {
		meta := sandbox.Metadata(sbClient.ParseMetadataMap(info.Metadata))
		params.Metadata = &meta
	}

	fmt.Printf("Creating sandbox from template %s...\n", info.TemplateID)
	sb, _, err := client.CreateAndWait(ctx, params)
	if err != nil {
		sbClient.PrintError("create sandbox failed: %v", err)
		return
	}
	sbClient.PrintSuccess("Sandbox %s created, connecting...", sb.ID())

	// When create session ends, kill the sandbox
	defer func() {
		fmt.Printf("\nKilling sandbox %s...\n", sb.ID())
		if kErr := sb.Kill(context.Background()); kErr != nil {
			// Ignore 404 errors: sandbox may have already been terminated by timeout
			if !strings.Contains(kErr.Error(), "404") {
				sbClient.PrintWarn("kill sandbox failed: %v", kErr)
			}
		}
	}()

	runTerminalSession(ctx, sb)
}
