//go:build !windows

package operations

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func notifyTerminalResize(ctx context.Context, resizeEvents chan<- struct{}) {
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	go func() {
		defer signal.Stop(sigWinch)
		for {
			select {
			case <-sigWinch:
				select {
				case resizeEvents <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
