//go:build !windows

package operations

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func notifyTerminalResize(_ context.Context, sigWinch chan<- os.Signal) {
	signal.Notify(sigWinch, syscall.SIGWINCH)
}
