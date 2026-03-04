package sandbox

import (
	"encoding/json"
	"fmt"
	"strings"

	sdkSandbox "github.com/qiniu/go-sdk/v7/sandbox"
)

// Output format constants.
const (
	FormatPretty = "pretty"
	FormatJSON   = "json"
)

// Connect timeout constants (in seconds).
const (
	// ConnectTimeoutInteractive is the timeout for interactive PTY sessions.
	ConnectTimeoutInteractive int32 = 300
	// ConnectTimeoutCommand is the timeout for non-interactive operations (kill, logs, metrics).
	ConnectTimeoutCommand int32 = 10
)

// PrintJSON outputs data as formatted JSON to stdout.
func PrintJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: marshal JSON failed: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// ParseStates parses a comma-separated state string into a slice of SandboxState.
func ParseStates(stateStr string) []sdkSandbox.SandboxState {
	parts := strings.Split(stateStr, ",")
	states := make([]sdkSandbox.SandboxState, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s != "" {
			states = append(states, sdkSandbox.SandboxState(s))
		}
	}
	return states
}
