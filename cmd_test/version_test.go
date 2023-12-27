//go:build unit

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestVersion(t *testing.T) {
	result := test.RunCmd(t, "version")
	if !strings.Contains(result, "UNSTABLE") {
		t.Fatal("version")
	}
}
