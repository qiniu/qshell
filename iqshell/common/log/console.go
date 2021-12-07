package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

// brush is a color join function
type brush func(string) string

const adapterConsole = "qn_console"

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;44"), // Emergency          white
	newBrush("1;36"), // AlertF              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // ErrorF              red
	newBrush("1;33"), // WarningF            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("1;37"), // DebugF              Background blue
}

// consoleWriter implements LoggerInterface and writes messages to terminal.
type consoleWriter struct {
	Level    int  `json:"level"`
	Colorful bool `json:"color"` //this filed is useful only when system's terminal supports color
}

// NewConsole create ConsoleWriter returning as LoggerInterface.
func newConsole() logs.Logger {
	cw := &consoleWriter{
		Level:    logs.LevelDebug,
		Colorful: true,
	}
	return cw
}

// Init init console logger.
// jsonConfig like '{"level":LevelTrace}'.
func (c *consoleWriter) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(jsonConfig), c)
}

// WriteMsg write message in console.
func (c *consoleWriter) WriteMsg(when time.Time, msg string, level int) (err error) {
	if level > c.Level {
		return
	}
	// alert 去除标识
	msg = strings.Replace(msg, "[A]  ", "", 1)
	if c.Colorful {
		msg = colors[level](msg)
	}
	if level == logs.LevelError {
		_, err = fmt.Fprintln(os.Stderr, msg)
	} else {
		_, err = fmt.Println(msg)
	}
	return
}

// Destroy implementing method. empty.
func (c *consoleWriter) Destroy() {

}

// Flush implementing method. empty.
func (c *consoleWriter) Flush() {

}

func init() {
	logs.Register(adapterConsole, newConsole)
}
