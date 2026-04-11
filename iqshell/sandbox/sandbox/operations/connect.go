package operations

import (
	"context"

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
		sbClient.PrintError("sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutInteractive})
	if err != nil {
		sbClient.PrintError("connect to sandbox %s failed: %v", info.SandboxID, err)
		return
	}
	sbClient.PrintSuccess("Connected to sandbox %s", sb.ID())

	runTerminalSession(ctx, sb)
}
