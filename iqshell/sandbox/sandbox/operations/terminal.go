package operations

import (
	"context"
	"os"
	"sync"
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

type terminalSize struct {
	width  int
	height int
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

func detectResize(previous terminalSize, width, height int, err error) (terminalSize, bool) {
	if err != nil {
		return previous, false
	}
	current := terminalSize{width: width, height: height}
	if current == previous {
		return previous, false
	}
	return current, true
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
	resizeEvents := make(chan struct{}, 1)
	notifyTerminalResize(ptyCtx, resizeEvents)

	currentSize := terminalSize{width: width, height: height}
	startTerminalResizeMonitor(ptyCtx, resizeEvents, currentSize, func() (int, int, error) {
		return term.GetSize(int(os.Stdin.Fd()))
	}, func(w, h int) {
		sb.Pty().Resize(ptyCtx, pid, sandbox.PtySize{
			Cols: uint32(w),
			Rows: uint32(h),
		})
	})

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

func startTerminalResizeMonitor(
	ctx context.Context,
	resizeEvents <-chan struct{},
	initialSize terminalSize,
	getSize func() (int, int, error),
	resize func(width, height int),
) {
	go func() {
		size := initialSize
		for {
			select {
			case <-resizeEvents:
				w, h, sErr := getSize()
				next, changed := detectResize(size, w, h, sErr)
				if changed {
					size = next
					resize(w, h)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
