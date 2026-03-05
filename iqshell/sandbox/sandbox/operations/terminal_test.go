package operations

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBatchedWriter_BuffersAndFlushes(t *testing.T) {
	var mu sync.Mutex
	var received [][]byte

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bw := newBatchedWriter(ctx, 50*time.Millisecond, func(ctx context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		cp := make([]byte, len(data))
		copy(cp, data)
		received = append(received, cp)
		return nil
	})

	// Write multiple chunks rapidly — should be batched together
	bw.Write([]byte("hello"))
	bw.Write([]byte(" world"))

	// Wait for flush interval to fire
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	count := len(received)
	mu.Unlock()

	if count == 0 {
		t.Fatal("expected at least 1 flush, got 0")
	}

	// The data should be concatenated in one or at most two flushes
	mu.Lock()
	var total int
	for _, r := range received {
		total += len(r)
	}
	mu.Unlock()

	if total != len("hello world") {
		t.Errorf("total flushed bytes = %d, want %d", total, len("hello world"))
	}

	bw.stop()
}

func TestBatchedWriter_StopFlushesRemaining(t *testing.T) {
	var mu sync.Mutex
	var received [][]byte

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bw := newBatchedWriter(ctx, 1*time.Second, func(ctx context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		cp := make([]byte, len(data))
		copy(cp, data)
		received = append(received, cp)
		return nil
	})

	bw.Write([]byte("final"))
	// Stop immediately — should trigger final flush
	bw.stop()

	// Give a tiny bit of time for the goroutine to finish
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("stop should trigger final flush")
	}
	if string(received[0]) != "final" {
		t.Errorf("final flush data = %q, want \"final\"", string(received[0]))
	}
}

func TestBatchedWriter_EmptyBufferNoFlush(t *testing.T) {
	var mu sync.Mutex
	flushCount := 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bw := newBatchedWriter(ctx, 20*time.Millisecond, func(ctx context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		flushCount++
		return nil
	})

	// Don't write anything — flush should be no-op
	time.Sleep(60 * time.Millisecond)
	bw.stop()

	mu.Lock()
	defer mu.Unlock()
	if flushCount != 0 {
		t.Errorf("empty buffer should not trigger flush, got %d flushes", flushCount)
	}
}

func TestBatchedWriter_MultipleWritesBatched(t *testing.T) {
	var mu sync.Mutex
	var received [][]byte

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use a long interval so all writes accumulate before flush
	bw := newBatchedWriter(ctx, 200*time.Millisecond, func(ctx context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		cp := make([]byte, len(data))
		copy(cp, data)
		received = append(received, cp)
		return nil
	})

	// Write 10 chunks rapidly
	for i := 0; i < 10; i++ {
		bw.Write([]byte("x"))
	}

	// Wait for one flush cycle
	time.Sleep(300 * time.Millisecond)
	bw.stop()

	mu.Lock()
	defer mu.Unlock()

	// All 10 bytes should be sent in 1-2 flushes (not 10)
	if len(received) > 2 {
		t.Errorf("expected at most 2 flushes for 10 rapid writes, got %d", len(received))
	}

	var total int
	for _, r := range received {
		total += len(r)
	}
	if total != 10 {
		t.Errorf("total flushed bytes = %d, want 10", total)
	}
}
