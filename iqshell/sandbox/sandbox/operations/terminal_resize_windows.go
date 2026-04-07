//go:build windows

package operations

import (
	"context"
	"os"
	"time"
)

const windowsResizePollInterval = 200 * time.Millisecond

func notifyTerminalResize(ctx context.Context, sigWinch chan<- os.Signal) {
	go func() {
		ticker := time.NewTicker(windowsResizePollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case sigWinch <- os.Interrupt:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
