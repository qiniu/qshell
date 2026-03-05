package sandbox

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// === LogLevelBadge tests ===

func TestLogLevelBadge_KnownLevels(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		badge := LogLevelBadge(level)
		if !strings.Contains(badge, strings.ToUpper(level)) {
			t.Errorf("LogLevelBadge(%q) = %q, should contain %q", level, badge, strings.ToUpper(level))
		}
	}
}

func TestLogLevelBadge_CaseInsensitive(t *testing.T) {
	badge := LogLevelBadge("INFO")
	if !strings.Contains(badge, "INFO") {
		t.Errorf("LogLevelBadge(\"INFO\") = %q, should contain INFO", badge)
	}
}

func TestLogLevelBadge_Unknown(t *testing.T) {
	badge := LogLevelBadge("custom")
	if !strings.Contains(badge, "CUSTOM") {
		t.Errorf("LogLevelBadge(\"custom\") = %q, should contain CUSTOM", badge)
	}
}

// === NewTable tests ===

func TestNewTable(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTable(&buf)
	tw.Write([]byte("A\tB\tC\n"))
	tw.Write([]byte("long\tshort\tx\n"))
	tw.Flush()
	output := buf.String()
	if !strings.Contains(output, "A") || !strings.Contains(output, "long") {
		t.Errorf("NewTable output unexpected: %s", output)
	}
}

// === FormatTimestamp tests ===

func TestFormatTimestamp_Zero(t *testing.T) {
	if got := FormatTimestamp(time.Time{}); got != "-" {
		t.Errorf("FormatTimestamp(zero) = %q, want \"-\"", got)
	}
}

func TestFormatTimestamp_Valid(t *testing.T) {
	ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	got := FormatTimestamp(ts)
	if got != "2024-01-15T10:30:00Z" {
		t.Errorf("FormatTimestamp = %q, want \"2024-01-15T10:30:00Z\"", got)
	}
}

// === FormatBytes tests ===

func TestFormatBytes_ExactMiB(t *testing.T) {
	// 512 MiB = 512 * 1024 * 1024 bytes
	if got := FormatBytes(512 * 1024 * 1024); got != "512 MiB" {
		t.Errorf("FormatBytes(512MiB) = %q, want \"512 MiB\"", got)
	}
}

func TestFormatBytes_FractionalMiB(t *testing.T) {
	// 1.5 MiB
	if got := FormatBytes(int64(1.5 * 1024 * 1024)); got != "1.5 MiB" {
		t.Errorf("FormatBytes(1.5MiB) = %q, want \"1.5 MiB\"", got)
	}
}

func TestFormatBytes_Zero(t *testing.T) {
	if got := FormatBytes(0); got != "0 MiB" {
		t.Errorf("FormatBytes(0) = %q, want \"0 MiB\"", got)
	}
}

// === FormatBytesHuman tests ===

func TestFormatBytesHuman_Bytes(t *testing.T) {
	if got := FormatBytesHuman(500); got != "500 B" {
		t.Errorf("FormatBytesHuman(500) = %q, want \"500 B\"", got)
	}
}

func TestFormatBytesHuman_KiB(t *testing.T) {
	got := FormatBytesHuman(2048)
	if got != "2.0 KiB" {
		t.Errorf("FormatBytesHuman(2048) = %q, want \"2.0 KiB\"", got)
	}
}

func TestFormatBytesHuman_GiB(t *testing.T) {
	got := FormatBytesHuman(2 * 1024 * 1024 * 1024)
	if got != "2.0 GiB" {
		t.Errorf("FormatBytesHuman(2GiB) = %q, want \"2.0 GiB\"", got)
	}
}

// === FormatMetadata tests ===

func TestFormatMetadata_Nil(t *testing.T) {
	if got := FormatMetadata(nil); got != "-" {
		t.Errorf("FormatMetadata(nil) = %q, want \"-\"", got)
	}
}

func TestFormatMetadata_Empty(t *testing.T) {
	if got := FormatMetadata(map[string]string{}); got != "-" {
		t.Errorf("FormatMetadata(empty) = %q, want \"-\"", got)
	}
}

func TestFormatMetadata_Single(t *testing.T) {
	got := FormatMetadata(map[string]string{"key": "val"})
	if got != "key=val" {
		t.Errorf("FormatMetadata single = %q, want \"key=val\"", got)
	}
}

func TestFormatMetadata_Multiple(t *testing.T) {
	got := FormatMetadata(map[string]string{"a": "1", "b": "2"})
	// Order is non-deterministic, check both keys exist
	if !strings.Contains(got, "a=1") || !strings.Contains(got, "b=2") {
		t.Errorf("FormatMetadata multiple = %q, should contain a=1 and b=2", got)
	}
}

// === FormatOptionalString tests ===

func TestFormatOptionalString_Nil(t *testing.T) {
	if got := FormatOptionalString(nil); got != "-" {
		t.Errorf("FormatOptionalString(nil) = %q, want \"-\"", got)
	}
}

func TestFormatOptionalString_Empty(t *testing.T) {
	s := ""
	if got := FormatOptionalString(&s); got != "-" {
		t.Errorf("FormatOptionalString(\"\") = %q, want \"-\"", got)
	}
}

func TestFormatOptionalString_Value(t *testing.T) {
	s := "hello"
	if got := FormatOptionalString(&s); got != "hello" {
		t.Errorf("FormatOptionalString(\"hello\") = %q, want \"hello\"", got)
	}
}
