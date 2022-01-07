package test

import "testing"

func Test(t *testing.T, args ...string) string {
	result := ""
	NewTestFlow(args...).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(DefaultTestErrorHandler(t)).Run()
	return result
}

func DefaultTestErrorHandler(t *testing.T) func(line string) {
	return func(line string) {
		t.Fail()
	}
}
