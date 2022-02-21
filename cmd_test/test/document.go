package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestDocument(cmdName string, t *testing.T) {
	prefix := fmt.Sprintf("# 简介\n`%s`", cmdName)
	result, _ := RunCmdWithError(cmdName, DocumentOption)
	if !strings.HasPrefix(result, prefix) {
		t.Fail()
		t.Log("document test fail for cmd:" + cmdName)
	}
}
