//go:build windows

package operations

import (
	"context"
	"time"
)

const windowsResizePollInterval = 200 * time.Millisecond

func notifyTerminalResize(ctx context.Context, resizeEvents chan<- struct{}) {
	go func() {
		ticker := time.NewTicker(windowsResizePollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
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
