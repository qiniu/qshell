package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ConnectInfo holds parameters for connecting to a sandbox.
type ConnectInfo struct {
	SandboxID string
}

// Connect connects to an existing sandbox terminal.
// When the terminal session ends, the sandbox is kept running.
func Connect(info ConnectInfo) {
	if info.SandboxID == "" {
		fmt.Println("Error: sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx := context.Background()
	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: 300})
	if err != nil {
		fmt.Printf("Error: connect to sandbox %s failed: %v\n", info.SandboxID, err)
		return
	}
	fmt.Printf("Connected to sandbox %s\n", sb.ID())

	runTerminalSession(ctx, sb)
}
