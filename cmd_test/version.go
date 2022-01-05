package cmd

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	NewTestFlow("version").ResultHandler(func(line string) {
		if !strings.Contains(line, "UNSTABLE") {
			t.Fatal("version")
		}
	}).ErrorHandler(defaultTestErrorHandler(t)).Run()
}
