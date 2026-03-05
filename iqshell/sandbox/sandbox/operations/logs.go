package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// LogsInfo holds parameters for viewing sandbox logs.
type LogsInfo struct {
	SandboxID string
	Level     string // Log level filter: DEBUG, INFO, WARN, ERROR (default: INFO)
	Limit     int32
	Format    string // pretty or json
	Follow    bool   // Keep streaming logs until sandbox is closed
	Loggers   string // Comma-separated logger name prefixes to filter
}

// Logs retrieves and displays sandbox logs.
func Logs(info LogsInfo) {
	if info.SandboxID == "" {
		sbClient.PrintError("sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
	defer signal.Stop(sigCh)

	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutCommand})
	if err != nil {
		sbClient.PrintError("connect to sandbox %s failed: %v", info.SandboxID, err)
		return
	}

	// Default level to INFO (matches e2b CLI)
	level := info.Level
	if level == "" {
		level = "INFO"
	}

	// Parse logger filters
	loggerPrefixes := sbClient.ParseLoggers(info.Loggers)

	// Async sandbox-done monitoring (non-blocking, checks every 5s)
	sandboxDone := make(chan struct{})
	if info.Follow {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					running, _ := sb.IsRunning(ctx)
					if !running {
						close(sandboxDone)
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	var start *int64

	for {
		params := &sandbox.GetLogsParams{
			Start: start,
		}
		if info.Limit > 0 && start == nil {
			params.Limit = &info.Limit
		}

		logs, lErr := sb.GetLogs(ctx, params)
		if lErr != nil {
			sbClient.PrintError("get sandbox logs failed: %v", lErr)
			return
		}

		if info.Format == sbClient.FormatJSON {
			if !info.Follow {
				sbClient.PrintJSON(logs)
				return
			}
			// In follow+json mode, print each batch
			if len(logs.Logs) > 0 || len(logs.LogEntries) > 0 {
				sbClient.PrintJSON(logs)
			}
		} else {
			printLogEntries(logs, level, loggerPrefixes)
		}

		if !info.Follow {
			if info.Format != sbClient.FormatJSON && len(logs.Logs) == 0 && len(logs.LogEntries) == 0 {
				fmt.Println("No logs found")
			}
			return
		}

		// Update start timestamp for next poll
		if len(logs.Logs) > 0 {
			lastTs := logs.Logs[len(logs.Logs)-1].Timestamp.UnixMilli() + 1
			start = &lastTs
		} else if len(logs.LogEntries) > 0 {
			lastTs := logs.LogEntries[len(logs.LogEntries)-1].Timestamp.UnixMilli() + 1
			start = &lastTs
		}

		// Check if sandbox is done or context cancelled
		select {
		case <-sandboxDone:
			if info.Format != sbClient.FormatJSON {
				fmt.Println("\nStopped printing logs — sandbox is closed")
			}
			return
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(400 * time.Millisecond)
	}
}

// printLogEntries outputs log entries with level and logger filtering, and colored level badges.
func printLogEntries(logs *sandbox.SandboxLogs, level string, loggerPrefixes []string) {
	if len(logs.LogEntries) > 0 {
		for _, entry := range logs.LogEntries {
			if !sbClient.IsLogLevelIncluded(string(entry.Level), level) {
				continue
			}
			// Filter by logger if specified
			logger := entry.Fields["logger"]
			if len(loggerPrefixes) > 0 {
				if !sbClient.MatchesLoggerPrefix(logger, loggerPrefixes) {
					continue
				}
			}

			// Clean logger name
			cleanLogger := sbClient.CleanLoggerName(logger)

			// Strip internal fields, keep user-relevant ones
			userFields := sbClient.StripInternalFields(entry.Fields)

			// Build output line
			var parts []string
			parts = append(parts,
				fmt.Sprintf("[%s]", entry.Timestamp.Format(time.RFC3339)),
				sbClient.LogLevelBadge(string(entry.Level)),
			)
			if cleanLogger != "" {
				parts = append(parts, fmt.Sprintf("[%s]", cleanLogger))
			}
			parts = append(parts, entry.Message)

			// Append user fields if any
			if len(userFields) > 0 {
				fieldParts := make([]string, 0, len(userFields))
				for k, v := range userFields {
					if k == "logger" {
						continue
					}
					fieldParts = append(fieldParts, fmt.Sprintf("%s=%s", k, v))
				}
				if len(fieldParts) > 0 {
					parts = append(parts, strings.Join(fieldParts, " "))
				}
			}

			fmt.Println(strings.Join(parts, " "))
		}
	} else if len(logs.Logs) > 0 {
		for _, l := range logs.Logs {
			fmt.Printf("[%s] %s\n",
				l.Timestamp.Format(time.RFC3339),
				l.Line,
			)
		}
	}
}
