package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ExecInfo holds parameters for executing a command in a sandbox.
type ExecInfo struct {
	SandboxID  string
	Command    []string
	Background bool
	Cwd        string
	User       string
	Envs       []string // KEY=VALUE pairs
}

// Exec executes a command in a sandbox.
// In foreground mode, stdout/stderr are streamed in real time and the exit code is propagated.
// In background mode, the process PID is printed and the command returns immediately.
func Exec(info ExecInfo) {
	if info.SandboxID == "" {
		sbClient.PrintError("sandbox ID is required")
		return
	}
	if len(info.Command) == 0 {
		sbClient.PrintError("command is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	sb, err := client.Connect(ctx, info.SandboxID, sandbox.ConnectParams{Timeout: sbClient.ConnectTimeoutCommand})
	if err != nil {
		sbClient.PrintError("connect to sandbox %s failed: %v", info.SandboxID, err)
		return
	}

	// Build command string. The SDK command API accepts a shell command string,
	// so quote argv-style inputs before joining to preserve spaces and metacharacters.
	cmd := shellQuoteArgs(info.Command)

	// Build command options
	var opts []sandbox.CommandOption
	if info.Cwd != "" {
		opts = append(opts, sandbox.WithCwd(info.Cwd))
	}
	if info.User != "" {
		opts = append(opts, sandbox.WithCommandUser(info.User))
	}
	if len(info.Envs) > 0 {
		envMap := parseEnvPairs(info.Envs)
		if len(envMap) > 0 {
			opts = append(opts, sandbox.WithEnvs(envMap))
		}
	}

	// Enable stdin forwarding if input is piped
	if isPipedStdin() {
		opts = append(opts, sandbox.WithStdin())
	}

	if info.Background {
		execBackground(ctx, sb, cmd, opts)
	} else {
		execForeground(ctx, sb, cmd, opts)
	}
}

// execForeground runs a command with real-time streaming output and propagates the exit code.
func execForeground(ctx context.Context, sb *sandbox.Sandbox, cmd string, opts []sandbox.CommandOption) {
	opts = append(opts,
		sandbox.WithOnStdout(func(data []byte) { os.Stdout.Write(data) }),
		sandbox.WithOnStderr(func(data []byte) { os.Stderr.Write(data) }),
	)

	handle, err := sb.Commands().Start(ctx, cmd, opts...)
	if err != nil {
		sbClient.PrintError("exec failed: %v", err)
		os.Exit(1)
	}

	// Wait for PID before setting up signal handling
	pid, err := handle.WaitPID(ctx)
	if err != nil {
		sbClient.PrintError("waiting for process start: %v", err)
		os.Exit(1)
	}

	// Forward SIGINT/SIGTERM to the remote process
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			_ = sb.Commands().Kill(ctx, pid)
		}
	}()

	// Forward piped stdin to the sandbox process
	if isPipedStdin() {
		go sendStdinToSandbox(ctx, sb, pid)
	}

	result, err := handle.Wait()
	signal.Stop(sigCh)
	if err != nil {
		sbClient.PrintError("exec failed: %v", err)
		os.Exit(1)
	}
	if result.Error != "" {
		sbClient.PrintError("%s", result.Error)
	}
	os.Exit(result.ExitCode)
}

// execBackground starts a command in the background and prints the PID.
func execBackground(ctx context.Context, sb *sandbox.Sandbox, cmd string, opts []sandbox.CommandOption) {
	handle, err := sb.Commands().Start(ctx, cmd, opts...)
	if err != nil {
		sbClient.PrintError("exec failed: %v", err)
		return
	}

	pid, err := handle.WaitPID(ctx)
	if err != nil {
		sbClient.PrintError("waiting for process start: %v", err)
		return
	}
	fmt.Printf("PID: %d\n", pid)

	// Background mode: send stdin data then return immediately
	if isPipedStdin() {
		sendStdinToSandbox(ctx, sb, pid)
	}
}

const stdinChunkSize = 64 * 1024 // 64KB, consistent with E2B CLI

// isPipedStdin checks if stdin is a pipe or file redirection (non-interactive terminal).
func isPipedStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

// sendStdinToSandbox reads from local stdin and forwards to the sandbox process.
// After reading all data, it calls CloseStdin to signal EOF.
func sendStdinToSandbox(ctx context.Context, sb *sandbox.Sandbox, pid uint32) {
	buf := make([]byte, stdinChunkSize)
	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			if sendErr := sb.Commands().SendStdin(ctx, pid, buf[:n]); sendErr != nil {
				return
			}
		}
		if err != nil {
			break
		}
	}
	// Signal EOF, ignore errors (process may have exited, or server may not support CloseStdin yet)
	_ = sb.Commands().CloseStdin(ctx, pid)
}

// parseEnvPairs parses KEY=VALUE pairs into a map.
func parseEnvPairs(pairs []string) map[string]string {
	envs := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && kv[0] != "" {
			envs[kv[0]] = kv[1]
		}
	}
	return envs
}

func shellQuoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, shellQuoteArg(arg))
	}
	return strings.Join(quoted, " ")
}

func shellQuoteArg(arg string) string {
	if arg == "" {
		return "''"
	}
	for _, r := range arg {
		if !isShellSafeRune(r) {
			return "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
		}
	}
	return arg
}

func isShellSafeRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= '0' && r <= '9':
		return true
	}
	switch r {
	case '_', '-', '.', '/', ':', '=', '@', '%', '+', ',':
		return true
	default:
		return false
	}
}
