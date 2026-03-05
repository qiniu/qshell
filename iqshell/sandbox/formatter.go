package sandbox

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// Color variables for consistent styling across the CLI.
var (
	ColorError   = color.New(color.FgRed)
	ColorSuccess = color.New(color.FgGreen)
	ColorWarn    = color.New(color.FgYellow)
	ColorInfo    = color.New(color.FgCyan)
	ColorMuted   = color.New(color.FgHiBlack)
)

// PrintError prints a red "Error: " prefixed message to stderr.
func PrintError(format string, a ...any) {
	ColorError.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
}

// PrintSuccess prints a green message to stdout.
func PrintSuccess(format string, a ...any) {
	ColorSuccess.Printf(format+"\n", a...)
}

// PrintWarn prints a yellow "Warning: " prefixed message to stderr.
func PrintWarn(format string, a ...any) {
	ColorWarn.Fprintf(os.Stderr, "Warning: "+format+"\n", a...)
}

// logLevelColors maps log levels to their badge styles (matching e2b CLI).
var logLevelColors = map[string]*color.Color{
	"debug": color.New(color.FgWhite),
	"info":  color.New(color.FgGreen),
	"warn":  color.New(color.FgYellow),
	"error": color.New(color.FgRed),
}

// LogLevelBadge returns a colorized log level label.
func LogLevelBadge(level string) string {
	lower := strings.ToLower(level)
	c, ok := logLevelColors[lower]
	if !ok {
		return fmt.Sprintf("%-5s", strings.ToUpper(level))
	}
	return c.Sprintf("%-5s", strings.ToUpper(level))
}

// NewTable creates a tabwriter.Writer for tab-separated column output.
func NewTable(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
}

// FormatTimestamp formats a time to RFC3339 or "-" for zero values.
func FormatTimestamp(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(time.RFC3339)
}

// FormatBytes formats bytes to "XXX MiB" (matching e2b output).
func FormatBytes(b int64) string {
	mib := float64(b) / (1024 * 1024)
	if mib == float64(int64(mib)) {
		return fmt.Sprintf("%d MiB", int64(mib))
	}
	return fmt.Sprintf("%.1f MiB", mib)
}

// FormatBytesHuman formats bytes using adaptive units (B, KiB, MiB, GiB, TiB).
func FormatBytesHuman(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTP"[exp])
}

// FormatMetadata formats a metadata map as "k1=v1, k2=v2" or "-".
func FormatMetadata(m map[string]string) string {
	if len(m) == 0 {
		return "-"
	}
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, k+"="+v)
	}
	return strings.Join(pairs, ", ")
}

// FormatOptionalString formats a *string to its value or "-".
func FormatOptionalString(s *string) string {
	if s == nil || *s == "" {
		return "-"
	}
	return *s
}
