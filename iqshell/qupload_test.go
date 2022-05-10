package iqshell

import (
	"strings"
	"testing"
)

func TestGetLineCount(t *testing.T) {
	lines := map[string]int64{
		"hello":            1,
		"\nhello":          2,
		"hello\nworld":     2,
		"hello\nworld\n":   2,
		"\nhello\nworld\n": 3,
		"\nhello\nworld":   3,
	}
	for key, want := range lines {
		got := GetLineCount(strings.NewReader(key))
		if got != want {
			t.Fatalf("key: %q, got=%d, want=%d\n", key, got, want)
		}
	}
}
