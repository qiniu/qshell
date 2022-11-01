package utils

import (
	"strings"
	"testing"
)

func TestJsonDataTrimComments01(t *testing.T) {
	content := `
## commit 01
aaa ## commit 02
bbb
## commit 03
`
	contentNew := string(JsonDataTrimComments([]byte(content)))
	t.Log(contentNew)

	if strings.Contains(contentNew, "#") {
		t.Fatal("result content #")
	}

	if !strings.Contains(contentNew, "aaa") {
		t.Fatal("result has no content aaa")
	}

	if !strings.Contains(contentNew, "bbb") {
		t.Fatal("result has no content bbb")
	}
}

func TestJsonDataTrimComments02(t *testing.T) {
	content := `
// commit 01
aaa // commit 02
bbb
// commit 03
`
	contentNew := string(JsonDataTrimComments([]byte(content)))
	t.Log(contentNew)

	if strings.Contains(contentNew, "//") {
		t.Fatal("result content //")
	}

	if !strings.Contains(contentNew, "aaa") {
		t.Fatal("result has no content aaa")
	}

	if !strings.Contains(contentNew, "bbb") {
		t.Fatal("result has no content bbb")
	}
}

func TestJsonDataTrimComments03(t *testing.T) {
	content := `
// commit 01
"aaa" # commit 02
bbb
// commit 03
`
	contentNew := string(JsonDataTrimComments([]byte(content)))
	t.Log(contentNew)

	if strings.Contains(contentNew, "#") {
		t.Fatal("result content #")
	}

	if strings.Contains(contentNew, "//") {
		t.Fatal("result content //")
	}

	if !strings.Contains(contentNew, "aaa") {
		t.Fatal("result has no content aaa")
	}

	if !strings.Contains(contentNew, "bbb") {
		t.Fatal("result has no content bbb")
	}
}
