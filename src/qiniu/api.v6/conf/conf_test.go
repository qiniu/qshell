package conf

import (
	"strings"
	"testing"
)

func TestUA(t *testing.T) {
	err := SetUser("")
	if err != nil {
		t.Fatal("expect no error")
	}
	err = SetUser("错误的UA")
	if err == nil {
		t.Fatal("expect an invalid ua format")
	}
	err = SetUser("Test0-_.")
	if err != nil {
		t.Fatal("expect no error")
	}
}

func TestFormat(t *testing.T) {
	str := "tesT0.-_"
	v := formatUserAgent(str)
	if !strings.Contains(v, str) {
		t.Fatal("should include user")
	}
	if !strings.HasPrefix(v, "QiniuGo/"+version) {
		t.Fatal("invalid format")
	}
}
