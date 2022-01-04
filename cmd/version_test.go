package cmd

import (
	"testing"
)

func TestVersion(t *testing.T) {
	setOsArgsAndRun([]string{"qshell", "version"})
}
