package operations

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// batchedWriter accumulates stdin data and flushes at regular intervals,
// reducing the number of SendInput calls (matching e2b CLI's BatchedQueue).
type batchedWriter struct {
	mu     sync.Mutex
	buf    []byte
	sendFn func(ctx context.Context, data []byte) error
	ctx    context.Context
	done   chan struct{}
}

// newBatchedWriter creates a batchedWriter that flushes at the given interval.
func newBatchedWriter(ctx context.Context, interval time.Duration, sendFn func(ctx context.Context, data []byte) error) *batchedWriter {
	bw := &batchedWriter{
		sendFn: sendFn,
		ctx:    ctx,
		done:   make(chan struct{}),
	}
	go bw.flushLoop(interval)
	return bw
}

// Write appends data to the buffer (called from stdin reader goroutine).
func (bw *batchedWriter) Write(data []byte) {
	bw.mu.Lock()
	bw.buf = append(bw.buf, data...)
	bw.mu.Unlock()
}

// flush sends any buffered data.
func (bw *batchedWriter) flush() {
	bw.mu.Lock()
	if len(bw.buf) == 0 {
		bw.mu.Unlock()
		return
	}
	data := bw.buf
	bw.buf = nil
	bw.mu.Unlock()

	bw.sendFn(bw.ctx, data)
}

// flushLoop periodically flushes the buffer until done.
func (bw *batchedWriter) flushLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bw.flush()
		case <-bw.done:
			bw.flush() // final flush
			return
		}
	}
}

// stop signals the flush loop to perform a final flush and exit.
func (bw *batchedWriter) stop() {
	close(bw.done)
}

// runTerminalSession creates a PTY session and handles stdin/stdout bridging.
func runTerminalSession(ctx context.Context, sb *sandbox.Sandbox) {
	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height = 80, 24
	}

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		sbClient.PrintError("failed to set raw mode: %v", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	ptyCtx, ptyCancel := context.WithCancel(ctx)
	defer ptyCancel()

	// Create PTY session
	handle, err := sb.Pty().Create(ptyCtx, sandbox.PtySize{
		Cols: uint32(width),
		Rows: uint32(height),
	}, sandbox.WithOnPtyData(func(data []byte) {
		os.Stdout.Write(data)
	}))
	if err != nil {
		sbClient.PrintError("create PTY failed: %v", err)
		return
	}

	// Wait for PID
	pid, err := handle.WaitPID(ptyCtx)
	if err != nil {
		sbClient.PrintError("wait for PTY PID failed: %v", err)
		return
	}

	// Handle terminal resize
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	defer signal.Stop(sigWinch)

	go func() {
		for {
			select {
			case <-sigWinch:
				w, h, sErr := term.GetSize(int(os.Stdin.Fd()))
				if sErr == nil {
					sb.Pty().Resize(ptyCtx, pid, sandbox.PtySize{
						Cols: uint32(w),
						Rows: uint32(h),
					})
				}
			case <-ptyCtx.Done():
				return
			}
		}
	}()

	// Keep-alive: periodically refresh sandbox timeout.
	// Matches e2b CLI: setInterval 5s, setTimeout 30s.
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		// Refresh immediately on session start
		sb.SetTimeout(ptyCtx, 30*time.Second)
		for {
			select {
			case <-ticker.C:
				sb.SetTimeout(ptyCtx, 30*time.Second)
			case <-ptyCtx.Done():
				return
			}
		}
	}()

	// Forward stdin to PTY using batched writer (10ms interval, matching e2b CLI)
	writer := newBatchedWriter(ptyCtx, 10*time.Millisecond, func(ctx context.Context, data []byte) error {
		return sb.Pty().SendInput(ctx, pid, data)
	})
	defer writer.stop()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, rErr := os.Stdin.Read(buf)
			if rErr != nil {
				ptyCancel()
				return
			}
			if n > 0 {
				writer.Write(buf[:n])
			}
		}
	}()

	// Wait for PTY to exit
	handle.Wait()
}
