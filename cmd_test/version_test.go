package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	result := test.RunCmd(t, "version")
	if !strings.Contains(result, "UNSTABLE") {
		t.Fatal("version")
	}
}
