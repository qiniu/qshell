package sandbox

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/muesli/termenv"
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

// boxStyle is the lipgloss style for boxed messages.
var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1)

// PrintBox prints a message inside a rounded box.
func PrintBox(msg string) {
	fmt.Println(boxStyle.Render(msg))
}

// PrintSuccessBox prints a success message inside a green-bordered box.
func PrintSuccessBox(msg string) {
	style := boxStyle.BorderForeground(lipgloss.Color("2")) // green
	fmt.Println(style.Render(msg))
}

// PrintErrorBox prints an error message inside a red-bordered box.
func PrintErrorBox(msg string) {
	style := boxStyle.BorderForeground(lipgloss.Color("1")) // red
	fmt.Fprintln(os.Stderr, style.Render(msg))
}

// Hyperlink renders a clickable terminal hyperlink using OSC 8 escape sequences.
// Falls back to "text (url)" format when the terminal does not support hyperlinks.
func Hyperlink(url, text string) string {
	output := termenv.NewOutput(os.Stdout)
	if output.HasDarkBackground() || !output.EnvNoColor() {
		// Use OSC 8 hyperlink if terminal supports it
		return output.Hyperlink(url, text)
	}
	return fmt.Sprintf("%s (%s)", text, url)
}

// FormatCodeBlock returns a styled block for displaying code snippets in terminal output.
func FormatCodeBlock(code, language string) string {
	// Use lipgloss muted style for the code block frame
	codeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
	return codeStyle.Render(code)
}
