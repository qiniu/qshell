package cmd

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	s := "hel lo"
	want := s

	items := ParseLine(s, ",")
	if len(items) != 1 {
		t.Fatalf("expected 1 string, got 0\n")
	}
	if items[0] != want {
		t.Fatalf("want = %s, got = %s\n", want, items[0])
	}
}
