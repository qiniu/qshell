package iqshell

import (
	"testing"
)

func TestBytesToReadable(t *testing.T) {
	sizes := map[int64]string{
		512:           "512B",
		1024:          "1.00KB",
		2048:          "2.00KB",
		1048576:       "1.00MB",
		1073741824:    "1.00GB",
		2073741824:    "1.93GB",
		1099511627776: "1.00TB",
	}

	for size, want := range sizes {
		got := BytesToReadable(size)
		if got != want {
			t.Fatalf("size got=%s, want=%s", got, want)
		}
	}
}
