package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// CreateInfo holds parameters for creating a sandbox.
type CreateInfo struct {
	TemplateID string
	Timeout    int32
}

// Create creates a new sandbox and connects to its terminal.
// When the terminal session ends, the sandbox is killed.
func Create(info CreateInfo) {
	if info.TemplateID == "" {
		fmt.Println("Error: template ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx := context.Background()
	params := sandbox.CreateParams{
		TemplateID: info.TemplateID,
	}
	if info.Timeout > 0 {
		params.Timeout = &info.Timeout
	}

	fmt.Printf("Creating sandbox from template %s...\n", info.TemplateID)
	sb, _, err := client.CreateAndWait(ctx, params)
	if err != nil {
		fmt.Printf("Error: create sandbox failed: %v\n", err)
		return
	}
	fmt.Printf("Sandbox %s created, connecting...\n", sb.ID())

	// When create session ends, kill the sandbox
	defer func() {
		fmt.Printf("\nKilling sandbox %s...\n", sb.ID())
		if kErr := sb.Kill(context.Background()); kErr != nil {
			fmt.Printf("Warning: kill sandbox failed: %v\n", kErr)
		}
	}()

	runTerminalSession(ctx, sb)
}
