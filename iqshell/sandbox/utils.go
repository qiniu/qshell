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

// DefaultSandboxTimeout is the default sandbox timeout in seconds when creating a sandbox
// for interactive use. The server default is 300s (5 min) which is too short for terminal sessions.
const DefaultSandboxTimeout int32 = 3600 // 1 hour

// PrintJSON outputs data as formatted JSON to stdout.
func PrintJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: marshal JSON failed: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// DefaultState is the default sandbox state filter (matches e2b CLI behavior).
const DefaultState = "running"

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

// ParseMetadata parses comma-separated key=value pairs into URL query format.
// Input:  "key1=value1,key2=value2"
// Output: "key1=value1&key2=value2"
func ParseMetadata(raw string) string {
	if raw == "" {
		return ""
	}
	pairs := strings.Split(raw, ",")
	var parts []string
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && strings.TrimSpace(kv[0]) != "" && strings.TrimSpace(kv[1]) != "" {
			parts = append(parts, strings.TrimSpace(kv[0])+"="+strings.TrimSpace(kv[1]))
		}
	}
	return strings.Join(parts, "&")
}

// logLevelOrder maps log levels to numeric order for hierarchical filtering.
// Higher levels include all lower levels (e.g., INFO includes INFO, WARN, ERROR).
var logLevelOrder = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}

// IsLogLevelIncluded checks if a log entry level should be included
// given the minimum allowed level (hierarchical: INFO includes WARN and ERROR).
func IsLogLevelIncluded(entryLevel, minLevel string) bool {
	if minLevel == "" {
		return true
	}
	entryOrd, ok1 := logLevelOrder[strings.ToLower(entryLevel)]
	minOrd, ok2 := logLevelOrder[strings.ToLower(minLevel)]
	if !ok1 || !ok2 {
		return true
	}
	return entryOrd >= minOrd
}

// MatchesLoggerPrefix checks if a logger name matches any of the allowed prefixes.
func MatchesLoggerPrefix(logger string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(logger, prefix) {
			return true
		}
	}
	return false
}

// ParseLoggers parses a comma-separated logger string into a slice of prefix strings.
func ParseLoggers(raw string) []string {
	if raw == "" {
		return nil
	}
	var result []string
	for _, l := range strings.Split(raw, ",") {
		l = strings.TrimSpace(l)
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}
