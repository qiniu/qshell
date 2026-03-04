package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

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
		fmt.Printf("Error: failed to set raw mode: %v\n", err)
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
		fmt.Printf("Error: create PTY failed: %v\n", err)
		return
	}

	// Wait for PID
	pid, err := handle.WaitPID(ptyCtx)
	if err != nil {
		fmt.Printf("Error: wait for PTY PID failed: %v\n", err)
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

	// Keep-alive: periodically refresh sandbox timeout
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sb.SetTimeout(ptyCtx, 5*time.Minute)
			case <-ptyCtx.Done():
				return
			}
		}
	}()

	// Forward stdin to PTY
	go func() {
		buf := make([]byte, 1024)
		for {
			n, rErr := os.Stdin.Read(buf)
			if rErr != nil {
				ptyCancel()
				return
			}
			if n > 0 {
				if sErr := sb.Pty().SendInput(ptyCtx, pid, buf[:n]); sErr != nil {
					return
				}
			}
		}
	}()

	// Wait for PTY to exit
	handle.Wait()
}
