package test

import (
	"testing"
)

func RunCmd(t *testing.T, args ...string) string {
	result := ""
	NewTestFlow(args...).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(DefaultTestErrorHandler(t)).Run()
	return result
}

func RunCmdWithError(args ...string) (string, string) {
	result := ""
	err := ""
	NewTestFlow(args...).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		err += line
	}).Run()
	return result, err
}

func DefaultTestErrorHandler(t *testing.T) func(line string) {
	return func(line string) {
		t.Fail()
	}
}
